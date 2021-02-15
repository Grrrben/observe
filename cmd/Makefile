# Go parameters
GOBUILD=go build
GOTEST=go test
GOGET=go get
BINARY_NAME=observe

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v
fmt:
	gofmt -w ../.
test:
	$(GOTEST) -v ../
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
