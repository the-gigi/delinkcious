import json
import os
import subprocess
import time
from getpass import getpass
from itertools import chain
from multiprocessing import Process

platform = None


def guess_platform():
    """Guess which platform the cluster is running on
    - minikube
    - kind
    - k3s
    - EKS
    - GKE
    """
    cfg = json.loads(run('kubectl config view -o json', echo=False))
    name = cfg['clusters'][0]['name']
    if name.startswith('gke'):
        return 'gke'
    if name.startswith('minikube'):
        return 'minikube'
    if name.startswith('eksctl.io'):
        return 'eks'
    if name.startswith('kind'):
        return 'kind'
    if 'k3' in name:
        return 'k3s'

    raise RuntimeError('Unknown platform for cluster: ' + name)


def run(cmd, echo=True):
    output = subprocess.check_output(cmd.split()).decode('utf-8')
    if output and output[-1] == '\n':
        output = output[:-1]
    if echo:
        print(output)

    return output


def enable_minikube_addons():
    """ """
    addons = 'ingress heapster efk metrics-server'.split()
    for addon in addons:
        run('minikube addons enable ' + addon)


def install_helm():
    """ """
    run('kubectl apply -f helm_rbac.yaml')
    run('helm init --service-account tiller')


def install_metrics_server():
    run("""helm install 'stable/metrics-server 
             --name metrics-server
             --version 2.0.4
             --namespace monitoring""")


def install_nats():
    """ """
    run('kubectl apply -f https://github.com/nats-io/nats-operator/releases/download/latest/00-prereqs.yaml')
    run('kubectl apply -f https://github.com/nats-io/nats-operator/releases/download/latest/10-deployment.yaml')

    time.sleep(3)
    run('kubectl apply -f nats_cluster.yaml')


def install_argocd():
    """ """
    run('kubectl create namespace argocd')
    run('kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml')

    # Initial password is the ArgoCD server pod id
    get_pod_name = 'kubectl get pods -n argocd -l app.kubernetes.io/name=argocd-server -o name'
    argocd_podname = run(get_pod_name)
    while not argocd_podname:
        argocd_podname = run(get_pod_name)
        time.sleep(1)

    # Wait for pod to be running
    get_phase = f'kubectl get {argocd_podname} -n argocd -o jsonpath=' + "'{.status.phase}'"
    phase = run(get_phase)
    while phase != "'Running'":
        phase = run(get_phase)
        time.sleep(1)

    # Port-forward to access the Argo CD server locally on port 8080:
    argocd_port_forward()

    # Login and update password
    curr_password = argocd_podname.split('/')[1]
    if curr_password[-1] == '\n':
        curr_password = curr_password[:-1]
    run(f'argocd login :8080 --insecure --username admin --password {curr_password}')
    new_password = os.environ['ARGOCD_PASSWORD']
    run(f'argocd account update-password --current-password {curr_password} --new-password {new_password}')


def install_nuclio():
    # Install nuclio in its own namespace
    run('kubectl create namespace nuclio')

    run('kubectl apply -f https://raw.githubusercontent.com/nuclio/nuclio/master/hack/k8s/resources/nuclio-rbac.yaml')
    run('kubectl apply -f https://raw.githubusercontent.com/nuclio/nuclio/master/hack/k8s/resources/nuclio.yaml')

    ## Get nuctl CLI and create a symlink
    ver = '1.1.5'
    path = f'nuctl-{ver}-darwin-amd64'
    run(f'curl -LO https://github.com/nuclio/nuclio/releases/download/{ver}/{path}')
    run(f'sudo mv {path} /usr/local/bin/nuctl')
    run('chmod +x /usr/local/bin/nuctl')

    argocd_port_forward()

    # Create an image pull secret, so Nuclio can deploy functions to our cluster.
    dockerhub_password = os.environ.get('DOCKERHUB_PASSWORD', getpass('Enter Dockerhub password: '))
    run(f"""kubectl create secret docker-registry registry-credentials -n nuclio
               --docker-username g1g1
               --docker-password {dockerhub_password}
               --docker-server registry.hub.docker.com
               --docker-email the.gigi@gmail.com""")


def deploy_link_checker():
    """ Deploy the link checker nuclio function"""
    orig_dir = os.getcwd()
    os.chdir('../../fun/link_checker')
    registry = 'gcr.io' if platform == 'gke' else 'g1g1'
    run('nuctl deploy link-checker -n nuclio -p . --registry g1g1')
    os.chdir(orig_dir)


def argocd_login():
    host = 'localhost:8080'
    password = os.environ['ARGOCD_PASSWORD']
    cmd = f'argocd login {host} --insecure --username admin --password {password}'
    output = run(cmd)
    print(output)


def argocd_port_forward():
    port_forward = 'kubectl port-forward -n argocd svc/argocd-server 8080:443'
    p = Process(target=lambda: run(port_forward))
    p.start()
    time.sleep(3)


def get_apps(namespace):
    """ """
    output = run(f'argocd app list -o wide')
    keys = 'name project namespace path repo'.split()
    apps = []
    lines = output.split('\n')
    headers = [h.lower() for h in lines[0].split()]
    for line in lines[1:]:
        items = line.split()
        app = {k: v for k, v in zip(headers, items) if k in keys}
        if app:
            apps.append(app)
    return apps


def create_project(project, cluster, namespace, description, repo):
    """ """
    cmd = f'argocd proj create {project} --description {description} -d {cluster},{namespace} -s {repo}'
    output = run(cmd)
    print(output)

    # Add access to resources
    cmd = f'argocd proj allow-cluster-resource {project} "*" "*"'
    output = run(cmd)
    print(output)


def create_app(name, project, namespace, repo, path):
    """ """
    cmd = f"""argocd app create {name} --project {project} 
                                       --dest-server https://kubernetes.default.svc 
                                       --dest-namespace {namespace} 
                                       --repo {repo} 
                                       --path {path}"""
    output = run(cmd)
    print(output)


def sync_app(name):
    """ """
    cmd = f'argocd app sync {name}'
    output = run(cmd)
    print(output)


def deploy_delinkcious_services():
    argocd_login()
    project = 'default'
    ns = 'default'
    description = 'Delicious-like link management system'
    repo = 'https://github.com/the-gigi/delinkcious'
    # create_project(project, 'https://kubernetes.default.svc', ns, '', repo)
    for app in 'link social-graph user news api-gateway'.split():
        service = app.replace('-', '_') + '_service'
        create_app(app, project, ns, repo, f'svc/{service}/k8s')
        sync_app(app)


def install_prometheus():
    """Install prometheus from the Helm chart

    Don't mess with the operator
    """
    run('helm install --name prometheus stable/prometheus')


def install_jeager():
    base_url = 'https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy'
    run('kubectl create namespace observability')
    run(f'kubectl crete -f {base_url}/crds/jaegertracing_v1_jaeger_crd.yaml')
    run(f'kubectl crete -f {base_url}/service_account.yaml')
    run(f'kubectl crete -f {base_url}/role.yaml')
    run(f'kubectl crete -f {base_url}/role_binding.yaml')
    run(f'kubectl crete -f jeager_in_memory.yaml')


def main():
    """Check the active cluster's platform and install components accordingly"""
    common_components = (
        install_helm,
        install_nats,
        install_argocd,
        deploy_delinkcious_services,
    )

    platform_components = dict(
        minikube=(
            enable_minikube_addons,
            install_nuclio,
            # deploy_link_checker,
        ),
        kind=(
            install_nuclio,
            install_metrics_server,
            deploy_link_checker,
            install_prometheus
        ),
        k3s=(
            install_nuclio,
            install_metrics_server,
            #deploy_link_checker,
            install_prometheus
        ),
        gke=(
            install_nuclio,
            #deploy_link_checker,
            install_prometheus
        ),
        aws=(
            install_prometheus
        )
    )

    global platform
    platform = guess_platform()
    components = chain.from_iterable((common_components, platform_components[platform]))
    components = platform_components[platform]
    for install_component in components:
        install_component()


if __name__ == '__main__':
    main()
