
.PHONY: all
all: build

.PHONY: gofmt
gofmt:
	@echo "Running gofmt"
	gofmt -s -w `find . -path ./vendor -prune -o -type f -name '*.go' -print`

.PHONY: build
build: 
	go build -o ./build/numalign ./cmd/main.go

.PHONY: deps-update
deps-update:
	go mod tidy && \
	go mod vendor
