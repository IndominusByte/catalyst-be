#!/bin/bash

export COMPOSE_IGNORE_ORPHANS=True

export BACKEND_IMAGE=catalyst-go
export BACKEND_IMAGE_TAG=development
export BACKEND_CONTAINER=catalyst-go-development
export BACKEND_HOST=catalyst-go.service
export BACKEND_STAGE=development

docker build -t "$BACKEND_IMAGE:$BACKEND_IMAGE_TAG" .
docker-compose -f ./manifest/docker-compose.development.yaml up -d --build
