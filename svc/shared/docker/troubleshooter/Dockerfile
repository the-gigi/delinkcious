FROM python:3

RUN apt-get update -y
RUN apt-get install -y vim \
                       dnsutils \
                       apt-transport-https \
                       ca-certificates

RUN curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg
RUN echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | tee /etc/apt/sources.list.d/kubernetes.list

RUN apt-get update -y
RUN apt-get install -y kubectl

RUN pip install kubernetes \
                httpie     \
                ipython

CMD bash