.PHONY: dev
dev:
	air

.PHONY: build
build:
	go build -o main .

.PHONY: run
run:
	./main

.PHONY: docs
docs:
	swag init

.PHONY: dev-with-docs
dev-with-docs: docs dev 

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: init
init: tidy docs

.PHONY: delete-users
delete-users:
	go run tools/delete_users.go