#!/bin/bash

export COMPOSE_IGNORE_ORPHANS=True

# postgresql
export POSTGRESQL_IMAGE=catalyst-postgresql
export POSTGRESQL_IMAGE_TAG=development
export POSTGRESQL_CONTAINER=catalyst-postgresql-development
export POSTGRESQL_HOST=catalyst-postgresql.service
export POSTGRESQL_USER=catalystdev
export POSTGRESQL_PASSWORD=inisecret
export POSTGRESQL_DB=catalyst
export POSTGRESQL_TIME_ZONE=Asia/Kuala_Lumpur
docker build -t "$POSTGRESQL_IMAGE:$POSTGRESQL_IMAGE_TAG" -f ./manifest-docker/Dockerfile.postgresql ./manifest-docker

# pgadmin
export PGADMIN_IMAGE=catalyst-pgadmin
export PGADMIN_IMAGE_TAG=development
export PGADMIN_CONTAINER=catalyst-pgadmin-development
export PGADMIN_HOST=catalyst-pgadmin.service
export PGADMIN_EMAIL=admin@catalyst.com
export PGADMIN_PASSWORD=inisecret
docker build -t "$PGADMIN_IMAGE:$PGADMIN_IMAGE_TAG" -f ./manifest-docker/Dockerfile.pgadmin ./manifest-docker

docker-compose -f ./manifest/docker-compose.development.yaml up -d --build
