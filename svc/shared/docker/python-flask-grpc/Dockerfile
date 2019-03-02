FROM python:alpine
RUN apk add build-base
COPY requirements.txt /tmp
WORKDIR /tmp
RUN pip install -r requirements.txt
