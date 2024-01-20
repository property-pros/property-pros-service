#!/bin/sh

source .env
PORT="$LISTEN_PORT"
PORT2="$GATEWAY_PORT"
export UNIQUE_TAG=$(date +%s)

gcloud config set project vireo-401203

gcloud builds submit --config ./cicd/cloudbuild.yaml --substitutions=_UNIQUE_TAG_=$UNIQUE_TAG

# gcloud run deploy property-pros-service --no-traffic --port $CONTAINER_PORT --tag us-west1-docker.pkg.dev/vireo-401203/pp/pps:$UNIQUE_TAG --platform managed --region us-west1 --allow-unauthenticated
gcloud run deploy property-pros-service --use-http2 --env-vars-file ./cicd/.env.yaml --port $PORT --port $PORT2 --platform managed --image us-west1-docker.pkg.dev/vireo-401203/pp/pps:$UNIQUE_TAG --region us-west1 --allow-unauthenticated
echo “deployed to gcloud run”

gcloud artifacts packages list --repository=pp --location=us-west1
