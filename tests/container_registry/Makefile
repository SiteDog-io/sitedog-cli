DOCKER_REGISTRY=docker.io
DOCKER_USERNAME=myusername
IMAGE_NAME=my-app
VERSION=latest

.PHONY: build push deploy

build:
	docker build -t $(IMAGE_NAME):$(VERSION) .
	docker tag $(IMAGE_NAME):$(VERSION) $(DOCKER_REGISTRY)/$(DOCKER_USERNAME)/$(IMAGE_NAME):$(VERSION)

push: build
	docker login $(DOCKER_REGISTRY) -u $(DOCKER_USERNAME) -p $(DOCKER_PASSWORD)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_USERNAME)/$(IMAGE_NAME):$(VERSION)

deploy: push
	kubectl set image deployment/$(IMAGE_NAME) app=$(DOCKER_REGISTRY)/$(DOCKER_USERNAME)/$(IMAGE_NAME):$(VERSION)

clean:
	docker rmi $(IMAGE_NAME):$(VERSION) || true
	docker rmi $(DOCKER_REGISTRY)/$(DOCKER_USERNAME)/$(IMAGE_NAME):$(VERSION) || true