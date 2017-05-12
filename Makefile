REPO=benniekrijger
IMAGE=todo-service-go
VERSION=latest

all: publish

clean:
	docker rmi ${IMAGE} &>/dev/null || true

build: clean
	bin/build.sh ${REPO} ${IMAGE} ${VERSION}

publish: build
	docker login -u ${REPO}
	docker push ${REPO}/${IMAGE}:${VERSION}
