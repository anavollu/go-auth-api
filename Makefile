include .env
export COGNITO_CLIENT_ID := ${COGNITO_CLIENT_ID}
export COGNITO_USER_POOL_ID := ${COGNITO_USER_POOL_ID}
export AWS_SDK_LOAD_CONFIG := ${AWS_SDK_LOAD_CONFIG}
DOCKER_IMAGE := go-auth-ecr/$(APP_NAME):latest
AWS_ECR_PATH := $(AWS_ACCOUNT_ID).dkr.ecr.$(REGION).amazonaws.com
ECR_REPOSITORY_NAME := go-auth-ecr

.PHONY: run
run:
	go run .

.PHONY: build
build:
	docker build -t $(DOCKER_IMAGE) .

.PHONY: deploy
deploy:
	aws ecr get-login-password --region $(REGION) | docker login --username AWS --password-stdin $(AWS_ECR_PATH)
	docker tag $(DOCKER_IMAGE) $(AWS_ECR_PATH)/$(ECR_REPOSITORY_NAME)
	docker push $(AWS_ECR_PATH)/$(ECR_REPOSITORY_NAME)
