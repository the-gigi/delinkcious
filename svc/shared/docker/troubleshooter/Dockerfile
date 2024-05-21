FROM python:3

RUN apt-get update -y
RUN apt-get install -y vim \
                       dnsutils \
                       apt-transport-https \
                       ca-certificates \
                       netcat-openbsd \
                       redis \
                       snap

RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" \
    && chmod +x kubectl \
    && mv kubectl /usr/bin/kubectl

RUN pip install kubernetes \
                httpie     \
                ipython

CMD bash