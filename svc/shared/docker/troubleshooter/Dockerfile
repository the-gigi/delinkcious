# syntax=docker/dockerfile:1.3
FROM python:3

ARG ARCH

RUN apt-get update -y
RUN apt-get install -y vim \
                       dnsutils \
                       apt-transport-https \
                       ca-certificates \
                       netcat-openbsd \
                       redis

RUN KUBECTL_URL="https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/${ARCH}/kubectl" \
    && curl -LO $KUBECTL_URL \
    && chmod +x kubectl \
    && mv kubectl /usr/bin/kubectl

RUN pip install kubernetes \
                httpie     \
                ipython

CMD bash
