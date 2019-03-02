# Base image for Python + Flask + gRPC

Building the gRPC Python client takes a LONG time. Building this base image
lets services that need these capabilities like the api-gateway service to have a quick build
by just copying its files into the container.



