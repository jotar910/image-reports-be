#!/bin/sh
docker-compose exec kafka kafka-topics --bootstrap-server localhost:29092 --list

echo -e 'Creating kafka topics'
docker-compose exec kafka kafka-topics --bootstrap-server localhost:29092 --create --if-not-exists --replication-factor 1 --partitions 1 --topic report-created
docker-compose exec kafka kafka-topics --bootstrap-server localhost:29092 --create --if-not-exists --replication-factor 1 --partitions 1 --topic report-deleted
docker-compose exec kafka kafka-topics --bootstrap-server localhost:29092 --create --if-not-exists --replication-factor 1 --partitions 1 --topic image-processed
docker-compose exec kafka kafka-topics --bootstrap-server localhost:29092 --create --if-not-exists --replication-factor 1 --partitions 1 --topic image-stored

echo -e 'Successfully created the following topics:'
docker-compose exec kafka kafka-topics --bootstrap-server localhost:29092 --list