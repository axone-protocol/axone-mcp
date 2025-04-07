SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c
.EXPORT_ALL_VARIABLES:

GOPATH ?= $(shell go env GOPATH)
CURDIR := $(shell pwd)

# Versions
ALPINE_VERSION          ?= 3.20
DOCKER_IMAGE_TAG        ?= latest
GOFUMPT_VERSION         ?= v0.7.0
GOLANG_VERSION          ?= 1.23
GOLANGCI_LINT_VERSION   ?= v2.0.2
GOTHANKS_VERSION        ?= latest
MOCKGEN_VERSION         ?= v0.5.0
TPARSE_VERSION			?= v0.17.0

# Some colors (if supported)
define get_color
$(shell tput -Txterm $(1) $(2) 2>/dev/null || echo "")
endef

COLOR_CYAN   = $(call get_color,setaf,6)
COLOR_GREEN  = $(call get_color,setaf,2)
COLOR_RED    = $(call get_color,setaf,1)
COLOR_RESET  = $(call get_color,sgr0,)
COLOR_WHITE  = $(call get_color,setaf,7)
COLOR_YELLOW = $(call get_color,setaf,3)

# Application
BINARY_NAME   = axone-mcp
BINARY_AMD64  = $(BINARY_NAME).amd64
$(info 🐚 ${COLOR_GREEN}Fetching ${COLOR_CYAN}version${COLOR_RESET} from ${COLOR_YELLOW}git${COLOR_RESET}...)
VERSION       = $(shell cat version)
$(info 🐚 ${COLOR_GREEN}Fetching ${COLOR_CYAN}commit${COLOR_RESET} from ${COLOR_YELLOW}git${COLOR_RESET}...)
COMMIT        = $(shell git log -1 --format='%H')

# Directories
TARGET_DIR = ./target
TOOLS_DIR  = $(TARGET_DIR)/tools

GOFUMPT_BIN       = $(TOOLS_DIR)/gofumpt/$(GOFUMPT_VERSION)/gofumpt
GOLANGCI_LINT_BIN = $(TOOLS_DIR)/golangci-lint/$(GOLANGCI_LINT_VERSION)/golangci-lint
GOTHANKS_BIN      = $(TOOLS_DIR)/gothanks/$(GOTHANKS_VERSION)/gothanks
MOCKGEN_BIN       = $(TOOLS_DIR)/mockgen/$(MOCKGEN_VERSION)/mockgen
TPARSE_BIN		  = $(TOOLS_DIR)/tparse/$(TPARSE_VERSION)/tparse

# Build options
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))
whitespace := $(subst ,, )
comma      := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# Flags
LD_FLAGS  = \
	-X github.com/axone-protocol/axone-mcp/internal/version.Name=$(BINARY_NAME) \
	-X github.com/axone-protocol/axone-mcp/internal/version.Version=$(VERSION)  \
	-X github.com/axone-protocol/axone-mcp/internal/version.Commit=$(COMMIT)
BUILD_FLAGS := -ldflags '$(LD_FLAGS)'
GO_BUILD := CGO_ENABLED=0 go build $(BUILD_FLAGS)

# Environments
ENVIRONMENTS = \
	darwin-amd64 \
	darwin-arm64 \
	linux-amd64 \
	windows-amd64
ENVIRONMENTS_TARGETS = $(addprefix build-go-, $(ENVIRONMENTS))

default: help

.PHONY: check-deps
check-deps: ## Check for required external tools (docker, curl)
	@$(call echo_msg, 🛠, Checking, dependencies, ...)
	@command -v docker > /dev/null || { echo "Error: docker is not installed. Aborting." >&2; exit 1; }
	@command -v curl > /dev/null || { echo "Error: curl is not installed. Aborting." >&2; exit 1; }

.PHONY: deps
deps: ## Download Go module dependencies
	@$(call echo_msg, 📥, Downloading, dependencies, ...)
	@go mod download

.PHONY: tools
tools: $(GOLANGCI_LINT_BIN) $(GOTHANKS_BIN) $(GOFUMPT_BIN) $(MOCKGEN_BIN) $(TPARSE_BIN) ## Install necessary development tools

.PHONY: thanks
thanks: tools ## Thanks to the contributors
	@$(call echo_msg, 🙏, Running, gothanks, ...)
	@$(GOTHANKS_BIN) -y | grep -v "is already"

.PHONY: build
build: build-go ## Build the project

.PHONY: build-go
build-go: deps  ## Build executable for the current environment (default build)
	@$(call echo_msg, 🏗️, Building, ${BINARY_NAME}, into ${COLOR_YELLOW}${TARGET_DIR})
	@$(call build-go,"","",${TARGET_DIR}/${BINARY_NAME})

.PHONY: build-go-all
build-go-all: $(ENVIRONMENTS_TARGETS) ## Build executables for all available environments

.PHONY: test
test: test-go ## Run all the tests

.PHONY: test-go
test-go: deps tools ## Run tests for Go source code
	@$(call echo_msg, 🧪, Testing, project, ...)
	@go test -v -coverprofile ./target/coverage.txt ./... -json | $(TPARSE_BIN)

.PHONY: lint
lint: lint-go ## Lint files

.PHONY: lint-go
lint-go: tools ## Lint Go source code
	@$(call echo_msg, 🔍, Linting, Go code, ...)
	@$(GOLANGCI_LINT_BIN) run ./...

.PHONY: format
format: format-go ## Format files

.PHONY: format-go
format-go: tools ## Format Go source code
	@$(call echo_msg, 📐, Formatting, Go source code, ...)
	@$(GOFUMPT_BIN) -w -l .

## Mock:
.PHONY: mock
mock: tools ## Generate all the mocks (for tests)
	@$(call echo_msg, 🧱, Generating, mocks, ...)
	@$(MOCKGEN_BIN) -destination=internal/mcp/dataverse_mocks.go -package=mcp github.com/axone-protocol/axone-sdk/dataverse QueryClient

.PHONY: docker
docker: build ## Build Docker container
	@$(call echo_msg, 📦, Building, Docker container, ...)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_AMD64) .
	@docker build -t $(BINARY_NAME):$(DOCKER_IMAGE_TAG) .

.PHONY: clean
clean: clean-artifacts clean-tools ## Clean up

.PHONY: clean-artifacts
clean-artifacts: ## Clean up build artifacts
	@$(call echo_msg, 🧹, Cleaning, build artifacts, ...)
	@rm -rf $(TARGET_DIR)

.PHONY: clean-tools
clean-tools: ## Clean up tools
	@$(call echo_msg, 🧹, Cleaning, tools, ...)
	@rm -rf $(TOOLS_DIR)

.PHONY: help
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${COLOR_YELLOW}make${COLOR_RESET} ${COLOR_GREEN}<target>${COLOR_RESET}'
	@echo ''
	@echo 'Targets:'
	@$(foreach V,$(sort $(.VARIABLES)), \
		$(if $(filter-out environment% default automatic,$(origin $V)), \
			$(if $(filter TOOL_%,$V), \
				export $V="$($V)";))) \
	awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${COLOR_YELLOW}%-20s${COLOR_GREEN}%s${COLOR_RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${COLOR_CYAN}%s${COLOR_RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST) | envsubst

$(TOOLS_DIR):
	@mkdir -p $(TOOLS_DIR)

$(GOLANGCI_LINT_BIN): | $(TOOLS_DIR)
	@$(call echo_msg, 📦, Installing, golangci-lint, $(COLOR_YELLOW)$(GOLANGCI_LINT_VERSION)$(COLOR_RESET)...)
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | \
		sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)
	@mkdir -p $(dir $(GOLANGCI_LINT_BIN))
	@cp $(shell go env GOPATH)/bin/golangci-lint $(dir $(GOLANGCI_LINT_BIN))

$(GOTHANKS_BIN):
	@$(call echo_msg, 📦, Installing, gothanks, $(COLOR_YELLOW)$(GOTHANKS_VERSION)$(COLOR_RESET)...)
	@mkdir -p $(dir $(GOTHANKS_BIN))
	@GOBIN="$$(cd $(dir $(GOTHANKS_BIN)) && pwd)" go install github.com/psampaz/gothanks@$(GOTHANKS_VERSION)

$(GOFUMPT_BIN):
	@$(call echo_msg, 📦, Installing, gofumpt, $(COLOR_YELLOW)$(GOFUMPT_VERSION)$(COLOR_RESET)...)
	@mkdir -p $(dir $(GOFUMPT_BIN))
	@GOBIN="$$(cd $(dir $(GOFUMPT_BIN)) && pwd)" go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)

$(MOCKGEN_BIN):
	@$(call echo_msg, 📦, Installing, mockgen, $(COLOR_YELLOW)$(MOCKGEN_VERSION)$(COLOR_RESET)...)
	@mkdir -p $(dir $(MOCKGEN_BIN))
	@GOBIN="$$(cd $(dir $(MOCKGEN_BIN)) && pwd)" go install go.uber.org/mock/mockgen@$(MOCKGEN_VERSION)

$(TPARSE_BIN):
	@$(call echo_msg, 📦, Installing, tparse, $(COLOR_YELLOW)$(TPARSE_VERSION)$(COLOR_RESET)...)
	@mkdir -p $(dir $(TPARSE_BIN))
	@GOBIN="$$(cd $(dir $(TPARSE_BIN)) && pwd)" go install github.com/mfridman/tparse@$(TPARSE_VERSION)

$(ENVIRONMENTS_TARGETS):
	@GOOS=$(word 3, $(subst -, ,$@)); \
    GOARCH=$(word 4, $(subst -, ,$@)); \
    FOLDER=${TARGET_DIR}/$$GOOS/$$GOARCH; \
    if [ $$GOOS = "windows" ]; then \
        EXTENSION=".exe"; \
	else \
		EXTENSION=""; \
    fi; \
    FILENAME=$$FOLDER/${BINARY_NAME}$$EXTENSION; \
	$(call echo_msg, 🏗️, Building, ${BINARY_NAME}, for environment ${COLOR_YELLOW}$$GOOS ($$GOARCH)${COLOR_RESET} into ${COLOR_YELLOW}$$FOLDER) && \
	$(call build-go,$$GOOS,$$GOARCH,$$FILENAME)

# Prints a message with color
# $(call echo_msg, <emoji>, <action>, <object>, <context>)
define echo_msg
	echo "$(strip $(1)) ${COLOR_GREEN}$(strip $(2))${COLOR_RESET} ${COLOR_CYAN}$(strip $(3))${COLOR_RESET} $(strip $(4))${COLOR_RESET}"
endef

# Build go executable
# $(call build-go, <os>, <arch>, <filename>)
define build-go
	GOOS=$1 GOARCH=$2 $(GO_BUILD) -o $3 ${CMD_ROOT}
endef
