#!/usr/bin/make -f

# ------------------------------------------------------------------------------
# General settings / variables
# ------------------------------------------------------------------------------

# import config. you can change default config with `make cnf="config_special.env" build`
cnf ?= config.env
include $(cnf)
export $(shell sed 's/=.*//' $(cnf))

# import deploy config. You can change the default deploy config with `make cnf="deploy_special.env" release`
dpl ?= deploy.env
include $(dpl)
export $(shell sed 's/=.*//' $(dpl))

# get current working dir
CWD=$(shell cd "$( dirname $0 )" >/dev/null 2>&1 && pwd )

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[35m%-30s\033[0m %s\n",$$1,$$2}' $(MAKEFILE_LIST) | egrep -Ev 'prod|dev'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[33m%-30s\033[0m %s\n",$$1,$$2}' $(MAKEFILE_LIST) | grep dev
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n",$$1,$$2}' $(MAKEFILE_LIST) | grep prod

.DEFAULT_GOAL := help

# ------------------------------------------------------------------------------
# DOCKER TASKS (DEV)
# ------------------------------------------------------------------------------

dev-login: ## Performs docker login to DEV image repository
	docker login $(DOCKER_REPO_DEV)

dev-build: ## Build the container image
	docker build -t $(APP_NAME) \
				 -f "$(CWD)/Dockerfile" \
				    "$(CWD)" \
				 --tag "$(DOCKER_REPO_DEV)/logstash-custom/$(APP_NAME):$(VERSION)" \
				 --no-cache \
				 --compress

dev-run: ## Run container on port configured in `config.env`
	docker run  -it -d --rm \
				--env-file=./config.env \
				--name="$(APP_NAME)" \
				-p=$(HOST_PORT_1):$(CONTAINER_PORT_1) \
				-p=$(HOST_PORT_2):$(CONTAINER_PORT_2) \
				"$(DOCKER_REPO_DEV)/logstash-custom/$(APP_NAME):$(VERSION)" \
	&& docker ps

# Docker publish
dev-push: ## Publish the `{version}` as defined at `config.env`
	@echo 'publish $(VERSION) to $(DOCKER_REPO_DEV)'
	docker push "$(DOCKER_REPO_DEV)/logstash-custom/$(APP_NAME):$(VERSION)"

# ------------------------------------------------------------------------------
# DOCKER TASKS (PROD)
# ------------------------------------------------------------------------------
prod-login: ## Performs docker login to production image repository
	docker login $(DOCKER_REPO_PROD)

prod-build: ## Build the container image
	docker build -t $(APP_NAME) \
				 -f "$(CWD)/Dockerfile" \
				    "$(CWD)" \
				 --tag "$(DOCKER_REPO_PROD)/logstash-custom/$(APP_NAME):$(VERSION)" \
				 --no-cache \
				 --compress

prod-run: ## Run container on port configured in `config.env`
	docker run  -it -d --rm \
				--env-file=./config.env \
				--name="$(APP_NAME)" \
				-p=$(HOST_PORT_1):$(CONTAINER_PORT_1) \
				-p=$(HOST_PORT_2):$(CONTAINER_PORT_2) \
				"$(DOCKER_REPO_PROD)/logstash-custom/$(APP_NAME):$(VERSION)" \
	&& docker ps

# Docker publish
prod-push: ## Publish the `{version}` as defined at `config.env`
	@echo 'publish $(VERSION) to $(DOCKER_REPO_PROD)'
	docker push "$(DOCKER_REPO_PROD)/logstash-custom/$(APP_NAME):$(VERSION)"

# ------------------------------------------------------------------------------
# HELPERS
# ------------------------------------------------------------------------------
version: ## Output the tagged version as defined at `config.env`
	@echo $(VERSION)

stop: ## Stop and remove a running container
	docker stop $(APP_NAME) \
				&& docker rm $$(docker ps -a -q -f status=exited) 2>/dev/null \
				|| echo "No stopped containers found"
	docker ps -a

clean: ## Removes all local images of APP_NAME
	docker rmi $$(docker images -qa $(APP_NAME)) -f

purge: ## Removes all dangling images
	docker rmi $$(docker images -q -f dangling=true) || echo "Nothing found"

# ------------------------------------------------------------------------------
# EOF
# ------------------------------------------------------------------------------