GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
VERSION?=0.0.1
DATE = $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_REV = $(shell git rev-parse --short HEAD)


BINARY_NAME=service
DOCKER_REGISTRY?=chuyms07
IMAGE_NAME=${DOCKER_REGISTRY}/${BINARY_NAME}
IMAGE_TAG=${IMAGE_NAME}:${GIT_REV}
BINARY_PATH=cmd/${BINARY_NAME}

SERVICE_PORT?=3000
EXPORT_RESULT?=false # for CI please set EXPORT_RESULT to true

.PHONY: all test build vendor

all: help

default: build_docker_img ## build docker image by default

## Build:
build: ## Build your project and put the output binary in _out/
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOCMD) build -buildmode=pie -a -ldflags ' -w -s  -X "main.version=${VERSION}" -X "main.buildDateTime=${DATE}" -X "main.gitRev=${GIT_REV}"' -o _out/$(BINARY_NAME) ./${BINARY_PATH}

docker_build: #Build the docker image for this service
	docker build -t ${IMAGE_NAME} .

docker_run: #Runs a container of the image for this service (make docker_build must run first)
	docker run -d -p 8080:8080 ${IMAGE_NAME}

docker_release:
	docker tag ${IMAGE_NAME} ${IMAGE_NAME}:latest
	docker tag ${IMAGE_NAME} ${IMAGE_NAME}:${VERSION}
	docker push ${IMAGE_NAME}:latest
	docker push ${IMAGE_NAME}:${VERSION}

docker_stop: #stops all docker containers of this image regardless of the tag
	docker stop $$(docker ps | grep ${IMAGE_NAME} | awk '{print $$1}')
	docker rm $$(docker ps -a | grep ${IMAGE_NAME} | awk '{print $$1}')

docker_build_run: docker_build docker_run

## Clean
clean: ## Remove build related files
	rm -fr ./_out