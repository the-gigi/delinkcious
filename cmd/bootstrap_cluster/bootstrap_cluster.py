"""Bootstrap the Delinkcious Kubernetes cluster


Pre-requisites:

```
pip install sh
```

Usage:

Make sure you have a minikube cluster with version >= 1.20 up and running

```
minikube start
```

Then:

```
python ./bootstrap_cluster.py
```

After installation:

```
kubectl port-forward -n argocd svc/argocd-server 8080:443
```

Then browse to localhost:8080 to connect to the ArgoCD server

The username is `admin`
Use the following command to get the password:
```
kubectl get secret argocd-initial-admin-secret -n argocd -o jsonpath='{.data.password}' | base64 --decode
```
"""

import json
import os
import sh
import subprocess
import time
from getpass import getpass
from itertools import chain
from multiprocessing import Process

platform = None


def verify_requirements():
    """Make sure requirements are present

    - kubectl (1.20+)
    - helm3   (version 3)
    - argocd CLI
    """

    result = json.loads(run('kubectl version -o json'))
    minorClientVersion = int(result['clientVersion']['minor'])
    minorServerVersion = int(result['serverVersion']['minor'])
    if minorClientVersion < 20:
        raise RuntimeError('kubectl version >= 1.20 is required. See https://kubernetes.io/docs/tasks/tools/')
    if minorServerVersion < 20:
        raise RuntimeError('Kubernetes master version must be >= 1.20.')

    result = run('helm version')
    result = result.split('Version:"')[1]
    if not result.startswith('v3'):
        raise RuntimeError('Helm version 3 is required. See https://helm.sh/docs/intro/install')

    result = run('which argocd')
    if not result.contains('not found'):
        raise RuntimeError('argocd CLI is required. See https://github.com/argoproj/argo-cd/blob/master/docs/cli_installation.md')

def guess_platform():
    """Guess which platform the cluster is running on
    - minikube
    - kind
    - k3s
    - EKS
    - GKE
    - AKS
    """
    cfg = json.loads(run('kubectl config view -o json', echo=False))
    name = cfg['clusters'][0]['name']
    if name.startswith('gke'):
        return 'gke'
    if name.startswith('aks'):
        return 'aks'
    if name.startswith('minikube'):
        return 'minikube'
    if name.startswith('eksctl.io'):
        return 'eks'
    if name.startswith('kind'):
        return 'kind'
    if 'k3' in name:
        return 'k3s'

    raise RuntimeError('Unknown platform for cluster: ' + name)


def kg(args, output='json'):
        """Get all resources of the specific kind

        If namespace is not provided get all resources in all namespaces

        Return the result as a Python object (parsed JSON)
        """
        args = args.split() + ['-o', output]
        result = sh.kubectl.get(*args)
        decoded = str(result.stdout, 'utf-8')
        if output == 'json':
            result = json.loads(decoded)
        else:
            result = decoded
            if result.endswith('\n'):
                result = result[:-1]

        return result


def run(cmd, echo=True):
    output = subprocess.check_output(cmd.split()).decode('utf-8')
    if output and output[-1] == '\n':
        output = output[:-1]
    if echo:
        print(output)

    return output


def enable_minikube_addons():
    """ """
    addons = 'ingress efk metrics-server'.split()
    for addon in addons:
        run('minikube addons enable ' + addon)


def install_metrics_server():
    run('helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server/')
    run("""helm upgrade --install 
           metrics-server 
           metrics-server/metrics-server 
           --version 2.0.4
           --namespace monitoring""")


def install_nats():
    """ """
    run('helm repo add nats https://nats-io.github.io/k8s/helm/charts/')
    run('helm repo update')
    run('helm upgrade --install nats-server nats/nats')


def install_argocd():
    """ """
    try:
        result = kg('namespace argocd', output='name')
    except Exception as e:
        result = e.stderr.decode('utf-8')
    if result != 'namespace/argocd':
        run('kubectl create  namespace argocd')
        run('kubectl create -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml')

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
    # Port-forward to access the Argo CD server locally on port 8080:
    argocd_port_forward()

    cmd = 'argocd login --core'
    output = run(cmd)
    print(output)


def argocd_prot_forward_target():
    port_forward = 'kubectl port-forward -n argocd svc/argocd-server 8080:443'
    run(port_forward)


def argocd_port_forward():
    p = Process(target=argocd_prot_forward_target)
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
    apps = 'link social-graph user news api-gateway'.split()
    #apps = ['api-gateway']
    for app in apps:
        service = app.replace('-', '_') + '_service'
        create_app(app, project, ns, repo, f'svc/{service}/k8s')
        sync_app(app)


def install_prometheus():
    """Install prometheus from the Helm chart

    Don't mess with the operator
    """
    run('helm repo add prometheus-community https://prometheus-community.github.io/helm-charts')
    run('helm repo update')
    run('helm install --generate-name prometheus-community/kube-prometheus-stack')


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
        install_nats,
        install_argocd,
        deploy_delinkcious_services,
    )

    platform_components = dict(
        minikube=(
            enable_minikube_addons,
            install_nuclio,
            #deploy_link_checker,
        ),
        kind=(
            install_nuclio,
            install_metrics_server,
            deploy_link_checker,
            install_prometheus,
        ),
        k3s=(
            install_nuclio,
            # install_metrics_server,
            # deploy_link_checker,
            # install_prometheus,
        ),
        gke=(
            install_nuclio,
            # deploy_link_checker,
            install_prometheus,
        ),
        aks=(
            # install_nuclio,
            # deploy_link_checker,
            install_prometheus,
        ),
        aws=(
            install_prometheus,
        )
    )

    # verify_requirements()

    global platform
    platform = guess_platform()
    components = chain.from_iterable((common_components, platform_components[platform]))
    # components = platform_components[platform]
    for install_component in components:
        print(f'running {install_component.__name__}()...')
        install_component()


if __name__ == '__main__':
    main()
