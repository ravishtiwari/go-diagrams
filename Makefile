# Makefile for the go-diagrams project

# Default target executed when no specific target is provided
.PHONY: all
all: build

# Variables
BINARY_NAME=go-diagrams
EXAMPLES_DIR=_examples
PACKAGE_NAME=go-diagrams-src.tar.gz

# Targets

# Compile the Go code (useful for checking compilation without producing a binary in the current dir)
.PHONY: compile
compile:
	@echo "Compiling Go code..."
	go build ./...

# Validate the Go code using go vet
.PHONY: validate
validate:
	@echo "Validating Go code..."
	go vet ./...

# Build the Go project's CLI tool
# Assumes cmd/gen/main.go is the main entry point for the CLI tool.
.PHONY: build
build:
	@echo "Building the project's CLI tool (${BINARY_NAME})..."
	go build -o ${BINARY_NAME} cmd/gen/main.go

# Run Go tests
.PHONY: test
test:
	@echo "Running Go tests..."
	go test ./...

# Create a gzipped tar archive of the source code
.PHONY: package
package:
	@echo "Packaging source code into ${PACKAGE_NAME}..."
	tar --exclude='.git' --exclude='${BINARY_NAME}' --exclude='${PACKAGE_NAME}' -czvf ${PACKAGE_NAME} .

# Run examples
# Iterates through all main.go files in the _examples directory.
# If a taskfile.yml is present, it runs 'go run main.go' then 'task -t taskfile.yml all'.
# Otherwise, it just runs 'go run main.go'.
.PHONY: run-examples
run-examples:
	@echo "Running examples..."
	@find ${EXAMPLES_DIR} -name "main.go" -print0 | while IFS= read -r -d $$'\0' example_main_file; do \
		example_dir=$$(dirname "$$example_main_file"); \
		echo "Running example in $$example_dir..."; \
		(cd "$$example_dir" && go run main.go); \
		if [ -f "$$example_dir/taskfile.yml" ]; then \
			echo "Found taskfile.yml in $$example_dir, running tasks..."; \
			(cd "$$example_dir" && task -t taskfile.yml all); \
		fi \
	done

# Clean build artifacts and packaged files
.PHONY: clean
clean:
	@echo "Cleaning build artifacts and packaged files..."
	rm -f ${BINARY_NAME}
	rm -f ${PACKAGE_NAME}
	@# The following command attempts to remove directories created by examples.
	@# It's a best-effort and might need adjustment based on actual example outputs.
	@find ${EXAMPLES_DIR} -mindepth 1 -maxdepth 1 -type d -print0 | while IFS= read -r -d $$'\0' example_sub_dir; do \
		if [ -d "$$example_sub_dir/go-diagrams" ]; then \
			echo "Removing $$example_sub_dir/go-diagrams..."; \
			rm -rf "$$example_sub_dir/go-diagrams"; \
		fi \
		if [ -d "$$example_sub_dir/test" ]; then \
			echo "Removing $$example_sub_dir/test..."; \
			rm -rf "$$example_sub_dir/test"; \
		fi \
	done
	@# Also remove top-level test directories if they exist from direct test runs
	@if [ -d "./go-diagrams" ]; then \
		echo "Removing ./go-diagrams..."; \
		rm -rf "./go-diagrams"; \
	fi
	@if [ -d "./test" ]; then \
		echo "Removing ./test..."; \
		rm -rf "./test"; \
	fi

# Help target to display available targets
.PHONY: help
help:
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:' Makefile | awk -F: '{print "\t" $$1}' | sort
