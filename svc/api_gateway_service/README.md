# API gateway Service

This is the [API gateway](https://microservices.io/patterns/apigateway.html) of Delinkcious. It is the only externally visible service.
It communicates with the other microservices using [API composition](https://microservices.io/patterns/data/api-composition.html):
- link service
- user service
- social graph service
- news service


# Implementing Social login

The API gaeway uses Github OAuth. You must add it as a Github Oauth application. See [Building Github OAuth Applications](https://developer.github.com/apps/building-oauth-apps/)


# Exposing to the world via port forwarding

The simplest way to expose the API gateway to the world is by port-forwrding the Kubernetes service directly:

```
kubectl port-forward svc/api-gateway 5000:5000
```







