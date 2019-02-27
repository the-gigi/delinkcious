FROM python:alpine
MAINTAINER Gigi Sayfan "the.gigi@gmail.com"
RUN apk add build-base
COPY . /api_gateway_service
WORKDIR /api_gateway_service
RUN pip install -r requirements.txt
EXPOSE 5000
ENTRYPOINT python run.py
