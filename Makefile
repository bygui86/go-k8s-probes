
# VARIABLES
export GO111MODULE = on


# CONFIG
.PHONY: help print-variables
.DEFAULT_GOAL := help


# ACTIONS

## infra

start-postgres :		## Run PostgreSQL in a container
	docker run -d --rm --name postgres \
		-e POSTGRES_PASSWORD=supersecret \
		-p 5432:5432 \
		postgres

stop-postgres :		## Stop PostgreSQL container
	docker stop postgres


## application

build :		## Build application
	go build

test :		## Run tests (required a running PostgreSQL instance)
	go test ./...

run :		## Run application from source code
	godotenv -f local.env go run main.go


## containerisation

__check-container-tag :
	@[ "$(CONTAINER_TAG)" ] || ( echo "Missing container tag (CONTAINER_TAG), please define it and retry"; exit 1 )

docker-build : __check-container-tag		## Build container
	docker build . -t bygui86/go-k8s-probes:$(CONTAINER_TAG)

docker-push : __check-container-tag		## Push container to Docker hub
	docker push bygui86/go-metrics:$(CONTAINER_TAG)


## helpers

help :		## Help
	@echo ""
	@echo "*** \033[33mMakefile help\033[0m ***"
	@echo ""
	@echo "Targets list:"
	@grep -E '^[a-zA-Z_-]+ :.*?## .*$$' $(MAKEFILE_LIST) | sort -k 1,1 | awk 'BEGIN {FS = ":.*?## "}; {printf "\t\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo ""

print-variables :		## Print variables values
	@echo ""
	@echo "*** \033[33mMakefile variables\033[0m ***"
	@echo ""
	@echo "- - - makefile - - -"
	@echo "MAKE: $(MAKE)"
	@echo "MAKEFILES: $(MAKEFILES)"
	@echo "MAKEFILE_LIST: $(MAKEFILE_LIST)"
	@echo "- - -"
	@echo ""
