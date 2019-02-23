FROM python:3
MAINTAINER Gigi Sayfan "the.gigi@gmail.com"
COPY . /quote-service
WORKDIR /quote-service
RUN pip install -r requirements.txt
EXPOSE 8000
ENTRYPOINT hug -f app.py
