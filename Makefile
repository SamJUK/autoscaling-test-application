MAKEFLAGS += --always-make
.DEFAULT_GOAL := help

build: ## Build docker image
	docker buildx build --platform linux/amd64 -t samjuk/autoscaling-test-server:latest .

push: ## Deploy docker image
	docker push --platform linux/amd64 samjuk/autoscaling-test-server:latest

help:  ## Prints the help document
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m  %-30s\033[0m %s\n", $$1, $$2}'