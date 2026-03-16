# Makefile for libnvme-go

# Variables
IMAGE_NAME ?= libnvme-go-test
CONTAINER_FILE ?= tests/Containerfile
PODMAN ?= podman

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  test          - Run unit tests in container"
	@echo "  build-image   - Build the test container image"
	@echo "  test-local    - Run tests locally (requires libnvme-devel)"
	@echo "  clean         - Remove the test container image"
	@echo "  shell         - Start an interactive shell in the test container"

.PHONY: build-image
build-image:
	$(PODMAN) build -t $(IMAGE_NAME) -f $(CONTAINER_FILE) .

.PHONY: test
test: build-image
	$(PODMAN) run --rm -v $(PWD):/workspace:Z $(IMAGE_NAME)

.PHONY: test-local
test-local:
	go test -v ./...

.PHONY: clean
clean:
	$(PODMAN) rmi -f $(IMAGE_NAME) || true

.PHONY: shell
shell: build-image
	$(PODMAN) run --rm -it -v $(PWD):/workspace:Z $(IMAGE_NAME) /bin/bash
