#!/bin/bash
set -e

DOCKERHUB_USER="hiamchubbybear"
VERSION="${1:-latest}"

services=("apigateway" "consultationservice" "identityservice" "notificationservice" "profileservice")

for service in "${services[@]}"; do
    IMAGE="$DOCKERHUB_USER/$service:$VERSION"
    echo "Checking $IMAGE on Docker Hub..."

    if docker manifest inspect "$IMAGE" > /dev/null 2>&1; then
        echo "$IMAGE already exists."

        read -p "Do you want to rebuild $service? [y/N]: " choice
        case "$choice" in
            y|Y )
                echo "Rebuilding $service..."
                docker build -f ./$service/Dockerfile -t "$IMAGE" ..
                docker push "$IMAGE"
                ;;
            * )
                echo "Skipping $service."
                ;;
        esac
    else
        echo "$IMAGE does not exist. Building..."
        docker build -f ./$service/Dockerfile -t "$IMAGE" ..
        docker push "$IMAGE"
    fi
done

echo "All services processed!"
