IMAGE=g1g1/py-kube:0.4

echo $DOCKERHUB_PASSWORD | docker login -u $DOCKERHUB_USERNAME --password-stdin

docker build . -t $IMAGE
docker push $IMAGE
