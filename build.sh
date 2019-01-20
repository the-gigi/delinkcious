#!/bin/bash

set -eo pipefail

IMAGE_PREFIX='g1g1'
STABLE_TAG='0.2'

TAG="${STABLE_TAG}.${CIRCLE_BUILD_NUM}"
ROOT_DIR="$(pwd)"
SVC_DIR="${ROOT_DIR}/svc"
cd $SVC_DIR
docker login -u $DOCKERHUB_USERNAME -p $DOCKERHUB_PASSWORD
for svc in *; do
    cd "${SVC_DIR}/$svc"
    if [[ ! -f Dockerfile ]]; then
        continue
    fi
    UNTAGGED_IMAGE=$(echo "${IMAGE_PREFIX}/delinkcious-${svc}" | sed -e 's/_/-/g' -e 's/-service//g')
    STABLE_IMAGE="${UNTAGGED_IMAGE}:${STABLE_TAG}"
    IMAGE="${UNTAGGED_IMAGE}:${TAG}"
    echo "image: $IMAGE"
    echo "stable image: ${STABLE_IMAGE}"
    docker build -t "$IMAGE" .
    docker tag "${IMAGE}" "${STABLE_IMAGE}"
    docker push "${IMAGE}"
    docker push "${STABLE_IMAGE}"
done
cd $ROOT_DIR