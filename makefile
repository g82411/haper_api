IMAGE_TAG?=$(shell echo `git rev-parse --short HEAD`)
DOCKER_PW=$(shell echo `aws ecr get-login-password --region ap-northeast-1`)
ECR_REPO=128945220432.dkr.ecr.ap-northeast-1.amazonaws.com
IMAGE_NAME=passpark-api
SOCKET_IMAGE_NAME=passpark-websocket
WORKER_IMAGE_NAME=passpark-worker

.PHONY: build_api migrate_new migrate_up build_socket build_worker

build_api:
	@docker login -u AWS -p $(DOCKER_PW) $(ECR_REPO)
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) -f Dockerfile .
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(ECR_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)
	docker push $(ECR_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)

build_socket:
	@docker login -u AWS -p $(DOCKER_PW) $(ECR_REPO)
	docker build -t $(SOCKET_IMAGE_NAME):$(IMAGE_TAG) -f socket.dockerfile .
	docker tag $(SOCKET_IMAGE_NAME):$(IMAGE_TAG) $(ECR_REPO)/$(SOCKET_IMAGE_NAME):$(IMAGE_TAG)
	docker push $(ECR_REPO)/$(SOCKET_IMAGE_NAME):$(IMAGE_TAG)


build_worker:
	@docker login -u AWS -p $(DOCKER_PW) $(ECR_REPO)
	docker build -t $(WORKER_IMAGE_NAME):$(IMAGE_TAG) -f worker.dockerfile .
	docker tag $(WORKER_IMAGE_NAME):$(IMAGE_TAG) $(ECR_REPO)/$(WORKER_IMAGE_NAME):$(IMAGE_TAG)
	docker push $(ECR_REPO)/$(WORKER_IMAGE_NAME):$(IMAGE_TAG)

migrate_new:
	migrate create -ext sql -dir ./migrations -seq $(name)

migrate_up:
	migrate -path ./migrations -database $(DSN) -verbose up