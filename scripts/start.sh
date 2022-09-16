#!/bin/sh

# Build database drectory
mkdir -p db-data

# Build kafka container.
docker-compose up -d init-kafka;
docker container wait init-kafka;

# Create kafka topics.
sh scripts/kafka.sh;

# Build the database containers.
#docker-compose up -d init-databases;
#docker container wait init-databases;

# Build all the other containers.
docker-compose up -d;