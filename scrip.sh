#!/bin/bash
set -e

DOCKERHUB_USER="hiamchubbybear"
VERSION="${1:-latest}"

# Login Docker Hub
echo "Logging into Docker Hub..."
echo "dckr_pat_LTuO4LaeRlB2Q3JCxLMKRKdG-Rs" | docker login -u "$DOCKERHUB_USER" --password-stdin

# Danh sách services
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
                # Build from root context with Dockerfile path
                docker build -f ./services/$service/Dockerfile -t "$IMAGE" .
                docker push "$IMAGE"
                ;;
            * )
                echo "Skipping $service."
                ;;
        esac
    else
        echo "$IMAGE does not exist. Building..."
        docker build -f ./services/$service/Dockerfile -t "$IMAGE" .
        docker push "$IMAGE"
    fi
done

echo "All services processed!"
