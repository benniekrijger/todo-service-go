#!/usr/bin/env bash

REPO=$1
IMAGE=$2
VERSION=$3

docker run --rm -it \
	  -v "$GOPATH":/gopath \
	  -v "$(pwd)":/app \
	  -e "GOPATH=/gopath" \
	  -w /app golang sh \
	  -c "CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags=\"-s\" -o main"
docker build -t ${REPO}/${IMAGE} .
rm main