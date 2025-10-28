PROJECT			 := github.com/globalcyberalliance/domain-trust-go
GO				 := $(shell which go 2>/dev/null)
GOFIELDALIGNMENT := $(shell which betteralign 2>/dev/null)
GOFUMPT			 := $(shell which gofumpt 2>/dev/null)
GOLINTER		 := $(shell which golangci-lint 2>/dev/null)
GONILAWAY        := $(shell which nilaway 2>/dev/null)
GO_BENCH_FLAGS	 := -short -bench=. -benchmem
GO_BENCH		 := $(GO) test $(GO_BENCH_FLAGS)
GO_BUILD         := $(GO_PRIVATE) CGO_ENABLED=0 $(GO) build -ldflags "-s -w" -trimpath
GO_FORMAT		 := $(GOFUMPT) -w
GO_OPTIMIZE		 := $(GOFIELDALIGNMENT) -fix
GO_PRIVATE		 := GOPRIVATE=github.com/globalcyberalliance
GO_TEST			 := $(GO) test -v -short
GO_TIDY			 := $(GO) mod tidy
TARGETS          = bin/client

.PHONY: ui

all: check-dependencies prepare optimize $(TARGETS) clean

dev: prepare $(TARGETS)

bin/%: $(shell find . -name "*.go" -type f)
	@echo "Building $@..."
	@if [ "$(MAKECMDGOALS)" != "dev" ]; then \
        cd build && $(GO_BUILD) -o ../$@ $(PROJECT)/cmd/$*; \
    else \
        $(GO_BUILD) -o $@ $(PROJECT)/cmd/$*; \
    fi

check-dependencies:
	@echo "Checking dependencies..."
	@if [ -z "${GO}" ]; then \
		echo "Cannot find 'go' in your $$PATH"; \
		exit 1; \
	fi
	@if [ -z "${GOFIELDALIGNMENT}" ]; then \
		echo "Cannot find 'betteralign' in your $$PATH"; \
		exit 1; \
	fi

clean:
	@echo "Cleaning temporary build directory..."
	@rm -rf build

docker:
	@if [ -z "$(GITHUB_TOKEN)" ]; then \
		echo "GITHUB_TOKEN not set"; \
		exit 1; \
	fi; \
	docker build \
		--build-arg GITHUB_TOKEN=$(GITHUB_TOKEN) \
		-t ghcr.io/globalcyberalliance/domain-trust:dev .

format:
	@if [ -z "${GOFUMPT}" ]; then \
		echo "Cannot find 'gofumpt' in your $$PATH"; \
		exit 1; \
	fi
	@echo "Formatting code..."
	@$(GO_FORMAT) -w $(PWD)

lint:
	@if [ -z "${GOLINTER}" ]; then \
		echo "Cannot find 'golangci-lint' in your $$PATH"; \
		exit 1; \
	fi
	@echo "Running linter..."
	@$(GOLINTER) run ./...

lint-apply:
	@if [ -z "${GOLINTER}" ]; then \
		echo "Cannot find 'golangci-lint' in your $$PATH"; \
		exit 1; \
	fi
	@echo "Running linter with autofix..."
	@$(GOLINTER) run --fix ./...

nil:
	@if [ -z "${GONILAWAY}" ]; then \
		echo "Cannot find 'nilaway' in your $$PATH"; \
		exit 1; \
	fi
	@echo "Running nilaway..."
	@$(GONILAWAY) ./...

optimize:
	@echo "Creating temporary build directory..."
	@cp -r cmd model go.* *.go ./build/
	@echo "Optimizing struct field alignment..."
	@cd build && $(GO_OPTIMIZE) ./... > /dev/null 2>&1 || true

prepare:
	@echo "Cleaning previous builds..."
	@rm -rf bin build
	@mkdir -p bin build
	@$(GO_TIDY)

setup:
	@echo "Installing dependencies (not needed for actual compile)..."
	@$(GO) install github.com/dkorunic/betteralign/cmd/betteralign@latest
	@$(GO) install mvdan.cc/gofumpt@latest
	@$(GO) install go.uber.org/nilaway/cmd/nilaway@latest
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

test:
	@echo "Running tests..."
	@$(GO_TEST) ./...