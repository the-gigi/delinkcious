# Delinkcious

A delicious-like link management platform implemented using Go microservices


# Directory Structure

## pkg
The core logic is implemented by libraries in this directory

## svc

The microservices are in this directory. They use the excellent [gokit](https://gokit.io) microservice framework.


## cmd

Various utilities and one-time commands live here


# CI/CD

For CI check out the .circleci file and build.sh

See https://circleci.com/gh/the-gigi/delinkcious/tree/master for status

For CD type: `kubectl port-forward -n argocd svc/argocd-server 8080:443`

Then browse to: https://localhost:8080
