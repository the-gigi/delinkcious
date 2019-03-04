# News Service

The news service is a gRPC microservice. Its gRPC contract lives in http:// Whenever you change the gRPC contract at `service/pb/news.proto` run these commands in the the `service/pb` directory:

```
protoc news.proto --go_out=plugins=grpc:.
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. news.proto
```

The first command generates `news.pb.go` that the service imports.

The second command generates two Python modules for Python consumers of the service.


It uses a multi-stage [Dockerfile](Dockerfile) to generate a lean and mean image from SCRATCH that just includes the Go binary. The system has a CI/CD pipeline, but you can also build and deploy it yourself.


## Build Docker image

```
$ docker build . -t g1g1/delinkcious-news-manager:${VERSION}
```

## Push to Registry

This will go by default to DockerHub. Make sure you're logged in:

```
$ docker login
```

Then push your image:

```
$ docker push g1g1/delinkcious-news-manager:${VERSION}
```

## Deploy to active Kubernetes cluster

If you want to push to a local minikube make sure your kubectl is pointed to the right cluster and type:

```
$ kubectl create -f k8s
```

## Exposing the NewsManager service locally

```
kubectl port-forward svc/news-manager 6060:6060
```









