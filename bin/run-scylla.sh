#!/usr/bin/env bash

echo "Removing old scylla instances..."

docker ps -aq --filter name=scylla | xargs docker rm -f

echo "Running new scylla instance..."

docker run \
  --detach \
  --name scylla \
  --publish 9042:9042 \
  scylladb/scylla