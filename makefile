IMAGE_TAG?=$(shell echo `git rev-parse --short HEAD`)
DOCKER_PW=$(shell echo `aws ecr get-login-password --region ap-northeast-1`)
ECR_REPO=843456404290.dkr.ecr.ap-northeast-1.amazonaws.com
IMAGE_NAME=hyper_api_image

.PHONY: build

build:
	@docker login -u AWS -p $(DOCKER_PW) $(ECR_REPO)
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) -f Dockerfile .
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(ECR_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)
	docker push $(ECR_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)