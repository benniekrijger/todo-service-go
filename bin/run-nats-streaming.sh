#!/usr/bin/env bash

echo "Removing old nats streaming instances..."

docker ps -aq --filter name=nats-streaming | xargs docker rm -f

echo "Running new nats streaming instance..."

: ${HOST:=$(ipconfig getifaddr en0)}
: ${HOST:=$(ipconfig getifaddr en1)}
: ${HOST:=$(ipconfig getifaddr en2)}
: ${HOST:=$(ipconfig getifaddr en3)}
: ${HOST:=$(ipconfig getifaddr en4)}

# could also expose zookeeper, add argument: -p 2181:2181 \
docker run \
  --detach \
  --name nats-streaming \
  -p 4222:4222 \
  nats-streaming