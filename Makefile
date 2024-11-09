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