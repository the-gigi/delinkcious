FROM python:3

RUN apt-get update -y
RUN apt-get install -y vim \
                       dnsutils
RUN pip install kubernetes \
                httpie     \
                ipython

CMD bash