# Social Graph Service

The social graph microservice. It uses a multi-stage Dockerfile to generate a lean and mean image from SCRATCH thst just includes the Go binary.


## Build Docker image

Use the correct version of course... 0.2 is just an example
```
docker build . -t g1g1/delinkcious-social-graph:0.2
```

## Push to Registry

```
docker push g1g1/delinkcious-social-graph:0.2
```







