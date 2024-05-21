IMAGE=g1g1/py-kube:0.4

echo $DOCKERHUB_PASSWORD | docker login -u $DOCKERHUB_USERNAME --password-stdin

# Create and use a new builder instance
docker buildx create --use

# Build and push the amd64 image
docker build --platform linux/amd64 \
    --build-arg ARCH=amd64 \
    -t ${IMAGE}-amd64 --push .

# Build and push the arm64 image
docker build --platform linux/arm64 \
    --build-arg ARCH=arm64 \
    -t ${IMAGE}-arm64 --push .
