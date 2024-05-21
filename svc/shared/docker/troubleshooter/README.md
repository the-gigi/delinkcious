# Docker image for the troubleshooter pod

See https://github.com/the-gigi/delinkcious/blob/master/svc/shared/k8s/troubleshooter.yaml

To build the image make sure you defined the following environment variables:

```
- DOCKERHUB_USERNAME
- DOCKERHUB_PASSWORD
```

Then, type
```
./build.sh
```
