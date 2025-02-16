#!/bin/bash

echo " Finding all services..."
services=($(find . -mindepth 1 -maxdepth 1 -type d -exec test -f "{}/pom.xml" \; -print | sed 's|./||'))

echo " Found services: ${services[*]}"

echo " Cleaning and packaging all services..."
for service in "${services[@]}"; do
    echo " Building $service..."
    cd $service
    mvn clean package -DskipTests -Dmaven.test.skip=true
    cd ..
done

echo " Stopping and removing all containers..."
docker compose down

echo " Pruning Docker volumes..."
docker volume prune -f

echo " Building and running Docker containers..."
docker compose up --build -d

echo " Deployment completed!"
