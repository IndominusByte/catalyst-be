#!/bin/bash

export COMPOSE_IGNORE_ORPHANS=True

export BACKEND_IMAGE=catalyst-go
export BACKEND_IMAGE_TAG=production
export BACKEND_CONTAINER=catalyst-go-production
export BACKEND_HOST=catalyst-go.service
export BACKEND_STAGE=production

docker build -t "$BACKEND_IMAGE:$BACKEND_IMAGE_TAG" .
docker-compose -f ./manifest/docker-compose.production.yaml up -d --build

