# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run
GOGET=$(GOCMD) get
BINARY_NAME=YamlApiServer

.PHONY: all test build run clean

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v
test:
	$(GOTEST) -v ./... -count=1
run: build
	$(GORUN) .
clean:
	$(GOCMD) clean
	rm -f $(BINARY_NAME)
