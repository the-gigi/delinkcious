# Delinkcious

A delicious-like link management platform implemented using Go microservices and deployed on Kubernetes.

The book [Hands-on Microservices with Kubernetes](https://www.amazon.com/Hands-Microservices-Kubernetes-scalable-microservices/dp/1789805465) describes how it was built from scratch.

# Directory Structure

## pkg
The core logic is implemented by libraries in this directory

## svc

The microservices are in this directory. They use the excellent [go-kit](https://gokit.io) microservice framework.

## cmd

Various utilities and one-time commands live here

## pb

Stands for protobuf. Contains gRPC contracts and generated code.

## fun

Serverless functions (Nuclio)

# Unit testing

Go to Delinkcious root directory and type: `ginkgo -r`

You should see something like:

```
[1556557699] LinkChecker Suite - 2/2 specs •• SUCCESS! 1.57716233s PASS
[1556557699] LinkManager Suite - 8/8 specs 2019/04/29 10:08:30 DB host: localhost DB port: 5432
•••••••• SUCCESS! 95.435161ms PASS
[1556557699] NewsManager Suite - 1/1 specs • SUCCESS! 322.678µs PASS
[1556557699] SocialGraphManager Suite - 6/6 specs 2019/04/29 10:08:30 DB host: localhost DB port: 5432
•••••• SUCCESS! 402.274617ms PASS
[1556557699] UserManager Suite - 6/6 specs 2019/04/29 10:08:31 DB host: localhost DB port: 5432
•••••• SUCCESS! 396.859071ms PASS

Ginkgo ran 5 suites in 11.589104359s
Test Suite Passed
```

# CI/CD

For CI check out the .circleci file and build.sh

See https://circleci.com/gh/the-gigi/delinkcious/tree/master for status

For CD type: `kubectl port-forward -n argocd svc/argocd-server 8080:443`

Then browse to: https://localhost:8080
