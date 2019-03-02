IMAGE=g1g1/delinkcious-python-flask-grpc:0.1
docker build . -t $IMAGE
docker push $IMAGE
