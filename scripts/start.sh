#!/bin/sh

# Build database drectory
mkdir -p db-data

# Build kafka container.
docker-compose up -d init-kafka;
docker container wait init-kafka;

# Create kafka topics.
sh scripts/kafka.sh;

# Build all the other containers.
docker-compose up -d;