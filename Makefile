GO_APP = goat

$(GO_APP):
	$(MAKE) up

.PHONY: up
up: ## start the goat service
	docker-compose up

build:
	docker-compose build goat

rebuild:
	docker-compose build --no-cache goat

test:
	docker-compose run --rm -e ENV=test goat go test ./...

shell: ## get a shell inside the goat container
	docker-compose run --rm goat /bin/sh

migrate:
	docker-compose run --rm goat go run ./cmd/migrate

help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

