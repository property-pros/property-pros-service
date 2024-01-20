#!/bin/bash
source .env
# Generate a unique tag based on the current timestamp
UNIQUE_TAG=$(date +%s)

# Build the Docker container with that unique tag
docker build -t cloud.canister.io:5000/vireo/property-pros-service:$UNIQUE_TAG .

# Push the Docker container to the specified Canister Docker registry
docker push cloud.canister.io:5000/vireo/property-pros-service:$UNIQUE_TAG

# Output the unique tag
echo "Successfully pushed: cloud.canister.io:5000/vireo/property-pros-service:$UNIQUE_TAG"

gcloud auth configure-docker

docker pull canister.io/vireo/property-pros-service:latest
docker tag canister.io/vireo/property-pros-service:latest us-docker.pkg.dev/vireo-clients/container/property-pros-service:latest

docker push us-docker.pkg.dev/vireo-clients/container/property-pros-service:latest

gcloud run deploy property-pros-service --use-http2 --image us-docker.pkg.dev/vireo-clients/container/property-pros-service:latest
echo "deployed to gcloud run"