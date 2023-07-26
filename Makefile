
.PHONY: all
all: run

.PHONY: gofmt
gofmt:
	@echo "Running gofmt"
	gofmt -s -w `find . -path ./vendor -prune -o -type f -name '*.go' -print`

.PHONY: build
build: 
	go build -o ./cmd/main.go

.PHONY: run
run: 
	go run ./cmd/main.go

.PHONY: deps-update
deps-update:
	go mod tidy && \
	go mod vendor