#!/bin/bash

DOCKERHUB_USERNAME="barnabasdoc22"

images=$(docker images --filter=reference='ethio*' --format "{{.Repository}}:{{.Tag}}")

for image in $images; do
    imagename=$(echo $image | cut -d':' -f1 | cut -d'/' -f2-)

    dockerhub_image="${DOCKERHUB_USERNAME}/${imagename}:latest"

    echo "Tagging $image as $dockerhub_image"
    docker tag "$image" "$dockerhub_image"

    echo "Pushing $dockerhub_image"
    docker push "$dockerhub_image"
done
