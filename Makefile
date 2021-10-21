GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BUILD_OUT_DIR=bin


PROJECT_DIR=./
BIN_NAME=go143
build:
	$(GOBUILD) -v -o $(BUILD_OUT_DIR)/$(BIN_NAME) $(PROJECT_DIR)

test:
	$(GOTEST) -v ./...

lint:
	golangci-lint run

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_OUT_DIR)/$(BIN_NAME)

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -v -o $(BUILD_OUT_DIR)/$(BIN_NAME)-linux $(PROJECT_DIR)

