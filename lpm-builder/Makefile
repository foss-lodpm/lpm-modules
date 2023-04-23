fmt-check:
	find . -name '*.go' -type f -not -path "*.git*" | xargs gofmt -d -s
	find . -name '*.go' -type f -not -path "*.git*" | xargs goimports -d -e

lint-check:
	golangci-lint run

check:
	$(MAKE) fmt-check
	$(MAKE) lint-check
	go mod verify

clean:
	rm -rf build

BUILD_FLAGS := -buildmode 'c-shared'

build: go.sum
	go build -mod=readonly $(BUILD_FLAGS) -o build/liblpm_builder.so cmd/main/main.go

.PHONY: build install test check
.ONESHELL:
.SILENT: