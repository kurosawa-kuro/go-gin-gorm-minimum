.PHONY: dev
dev:
	air

.PHONY: build
build:
	go build -o main .

.PHONY: run
run:
	./main 