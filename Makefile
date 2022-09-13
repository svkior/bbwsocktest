.DEFAULT_GOAL := help

.PHONY: create-network
create-network: ## create a local docker network
	@docker network inspect local >/dev/null || docker network create local
# Main targets
.PHONY: build-base
build-base: create-network ## build base golang working image
	@docker build \
		-t localhost:5000/air:latest \
		-f  deploy/base-image/Dockerfile \
		deploy/base-image/

.PHONY: watch-wsclient
watch-wsclient: create-network## start wsclient in autoreload mode
	@docker run -it --rm --name wsclient \
		--network local \
		--env-file "./deploy/wsclient.watch.env" \
		-v ${PWD}:/project \
		-v golang-cache-vol:/go/pkg/mod \
		-v go-build-vol:/root/.cache/go-build \
		--workdir="/project" \
		localhost:5000/air -c deploy/wsclient.air.toml

.PHONY: help
help: ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'