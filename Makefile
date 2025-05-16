# Protocol Buffer Compilation Makefile

# Directories
PROTO_DIR := src/protobuf
OUT_DIR := src/protobuf

# Find all proto files
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

# Generate output directories based on proto files
PROTO_OUTPUTS := $(patsubst $(PROTO_DIR)/%.proto,$(OUT_DIR)/%,$(PROTO_FILES))

# Default target
all: compile-proto

# Install required protobuf tools
install-proto-tools:
	@echo "Installing Protocol Buffer compiler tools..."
	@if command -v apt-get >/dev/null 2>&1; then \
		echo "Detected Ubuntu/Debian system, installing protoc with apt..."; \
		sudo apt-get update && sudo apt-get install -y protobuf-compiler; \
	else \
		echo "Could not detect apt-get. Please install protobuf compiler manually:"; \
		echo "For Ubuntu/Debian: sudo apt-get install -y protobuf-compiler"; \
		echo "For other systems, see: https://grpc.io/docs/protoc-installation/"; \
		exit 1; \
	fi
	@echo "Checking Go installation..."
	@if command -v go >/dev/null 2>&1; then \
		echo "Go found in PATH, installing Go protobuf plugins..."; \
		go install google.golang.org/protobuf/cmd/protoc-gen-go@latest; \
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest; \
	elif [ -f /usr/local/go/bin/go ]; then \
		echo "Go found in /usr/local/go/bin, installing Go protobuf plugins..."; \
		export PATH=$$PATH:/usr/local/go/bin; \
		/usr/local/go/bin/go install google.golang.org/protobuf/cmd/protoc-gen-go@latest; \
		/usr/local/go/bin/go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest; \
	elif [ -f /usr/lib/go/bin/go ]; then \
		echo "Go found in /usr/lib/go/bin, installing Go protobuf plugins..."; \
		export PATH=$$PATH:/usr/lib/go/bin; \
		/usr/lib/go/bin/go install google.golang.org/protobuf/cmd/protoc-gen-go@latest; \
		/usr/lib/go/bin/go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest; \
	elif [ -f /usr/bin/go ]; then \
		echo "Go found in /usr/bin, installing Go protobuf plugins..."; \
		/usr/bin/go install google.golang.org/protobuf/cmd/protoc-gen-go@latest; \
		/usr/bin/go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest; \
	else \
		echo "Error: Go not found. Please install Go using:"; \
		echo "  sudo apt-get install golang-go"; \
		echo "And ensure Go is in your PATH."; \
		exit 1; \
	fi
	@echo "Checking protoc installation..."
	@which protoc || echo "Warning: protoc still not found. Installation may have failed."

# Compile all proto files
compile-proto: $(PROTO_OUTPUTS)

# Rule to compile each proto file
$(OUT_DIR)/%: $(PROTO_DIR)/%.proto
	@echo "Compiling $<..."
	@mkdir -p $@
	@PACKAGE_NAME=$$(grep -m 1 "package" $< | sed 's/package\s*\([^;]*\);/\1/' | tr '.' '/'); \
	BASENAME=$$(basename $< .proto); \
	if [ ! -z "$$PACKAGE_NAME" ]; then \
		mkdir -p $@/$$PACKAGE_NAME; \
		IMPORT_PATH="github.com/snappfood/$$BASENAME/$$PACKAGE_NAME"; \
	else \
		IMPORT_PATH="github.com/snappfood/$$BASENAME"; \
	fi; \
	TMP_FILE=$$(mktemp); \
	if grep -q "option go_package" $<; then \
		cp $< $$TMP_FILE; \
	else \
		echo "Adding go_package option..."; \
		awk -v import_path="$$IMPORT_PATH" 'NR==2 { print "option go_package = \"" import_path "\";\n" } { print }' $< > $$TMP_FILE; \
	fi; \
	protoc --proto_path=$$(dirname $$TMP_FILE) \
		--go_out=$@ --go_opt=paths=source_relative \
		--go-grpc_out=$@ --go-grpc_opt=paths=source_relative \
		$$TMP_FILE; \
	rm $$TMP_FILE
	@echo "Generated code for $< in $@"

# Clean generated files
clean-proto:
	@echo "Cleaning generated Protocol Buffer files..."
	@for dir in $(PROTO_OUTPUTS); do \
		if [ -d $$dir ]; then \
			echo "Removing $$dir"; \
			rm -rf $$dir; \
		fi; \
	done

# Help target
help:
	@echo "Protocol Buffer Compilation Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  all               : Default target, same as compile-proto"
	@echo "  compile-proto     : Compile all proto files"
	@echo "  clean-proto       : Remove all generated files"
	@echo "  install-proto-tools: Install required protobuf compiler tools"
	@echo "  help              : Show this help message"

.PHONY: all compile-proto clean-proto install-proto-tools help
