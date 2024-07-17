# Define project name
PROJECT_NAME := service

# Define directories
BIN_DIR := bin
SRC_DIRS := ./src/cmd  ./src/pkg

# Define output binary
BINARY := $(BIN_DIR)/$(PROJECT_NAME)

# Define Go command
GO := go

.PHONY: run

# Default target: build the project
all: run

# Build the project
build: $(BINARY)

# Build binary
$(BINARY): $(SRC_DIRS)
	 mkdir -p $(BIN_DIR)
	 $(GO) build -o $(BINARY)  ./src/cmd/main.go

sw:
	swag init -g cmd/main.go --parseDependency --parseInternal

run:
	set APP_ENV=dev
	air -c .airm.toml

# Run tests
test:
	$(GO) test -v ./...


# Clean build artifacts
clean:
	rm -rf $(BIN_DIR)


deploy:
	docker-compose -f deploy/docker-compose.yml up -d --build

down:
	docker-compose -f deploy/docker-compose.yml down
