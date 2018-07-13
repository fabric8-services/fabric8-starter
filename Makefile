PROJECT_NAME=fabric8-starter
PACKAGE_NAME:=github.com/fabric8-services/$(PROJECT_NAME)

CUR_DIR=$(shell pwd)
TMP_PATH=$(CUR_DIR)/tmp
INSTALL_PREFIX=$(CUR_DIR)/bin
VENDOR_DIR=vendor

# declares variable that are OS-sensitive
SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
ifeq ($(OS),Windows_NT)
include $(SELF_DIR)Makefile.win
else
include $(SELF_DIR)Makefile.lnx
endif

# -------------------------------------------------------------------
# help!
# -------------------------------------------------------------------

.PHONY: help
help: ## Print this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

# -------------------------------------------------------------------
# required tools
# -------------------------------------------------------------------

# Find all required tools:
GIT_BIN := $(shell command -v $(GIT_BIN_NAME) 2> /dev/null)
DEP_BIN_DIR := $(TMP_PATH)/bin
DEP_BIN := $(DEP_BIN_DIR)/$(DEP_BIN_NAME)
DEP_VERSION=v0.4.1
GO_BIN := $(shell command -v $(GO_BIN_NAME) 2> /dev/null)

$(INSTALL_PREFIX):
	mkdir -p $(INSTALL_PREFIX)
$(TMP_PATH):
	mkdir -p $(TMP_PATH)

.PHONY: prebuild-check
prebuild-check: $(TMP_PATH) $(INSTALL_PREFIX) 
# Check that all tools where found
ifndef GIT_BIN
	$(error The "$(GIT_BIN_NAME)" executable could not be found in your PATH)
endif
ifndef DEP_BIN
	$(error The "$(DEP_BIN_NAME)" executable could not be found in your PATH)
endif
ifndef GO_BIN
	$(error The "$(GO_BIN_NAME)" executable could not be found in your PATH)
endif

# -------------------------------------------------------------------
# deps
# -------------------------------------------------------------------
$(DEP_BIN_DIR):
	mkdir -p $(DEP_BIN_DIR)

.PHONY: deps 
deps: $(DEP_BIN) $(VENDOR_DIR) ## Download the build dependencies.

# install dep in a the tmp/bin dir of the repo
$(DEP_BIN): $(DEP_BIN_DIR) 
	@echo "Installing 'dep' $(DEP_VERSION) at '$(DEP_BIN_DIR)'..."
	mkdir -p $(DEP_BIN_DIR)
ifeq ($(UNAME_S),Darwin)
	@curl -L -s https://github.com/golang/dep/releases/download/$(DEP_VERSION)/dep-darwin-amd64 -o $(DEP_BIN) 
	@cd $(DEP_BIN_DIR) && \
	curl -L -s https://github.com/golang/dep/releases/download/$(DEP_VERSION)/dep-darwin-amd64.sha256 -o $(DEP_BIN_DIR)/dep-darwin-amd64.sha256 && \
	echo "1544afdd4d543574ef8eabed343d683f7211202a65380f8b32035d07ce0c45ef  dep" > dep-darwin-amd64.sha256 && \
	shasum -a 256 --check dep-darwin-amd64.sha256
else
	@curl -L -s https://github.com/golang/dep/releases/download/$(DEP_VERSION)/dep-linux-amd64 -o $(DEP_BIN)
	@cd $(DEP_BIN_DIR) && \
	echo "31144e465e52ffbc0035248a10ddea61a09bf28b00784fd3fdd9882c8cbb2315  dep" > dep-linux-amd64.sha256 && \
	sha256sum -c dep-linux-amd64.sha256
endif
	@chmod +x $(DEP_BIN)

$(VENDOR_DIR): Gopkg.toml
	@echo "checking dependencies with $(DEP_BIN_NAME)"
	@$(DEP_BIN) ensure -v 


.PHONY: build
build: deps generate-assets  ## Build the binary
ifeq ($(OS),Windows_NT)
	go build -v $(LDFLAGS) -o "$(shell cygpath --windows '$(BINARY_BIN)')"
else
	go build -v $(LDFLAGS) -o $(BINARY_BIN)
endif


BINDATA_BIN:=$(VENDOR_DIR)/github.com/go-bindata/go-bindata/go-bindata/$(BINDATA_NAME)

.PHONY: generate-assets
generate-assets: prebuild-check $(BINDATA_BIN) ## Generate the assets file
	$(BINDATA_BIN) -pkg bootstrap -o bootstrap/bindata.go assets/...

$(BINDATA_BIN): deps
	cd $(VENDOR_DIR)/github.com/go-bindata/go-bindata/go-bindata && \
	go build -o $(BINDATA_NAME) . && \
	chmod u+x $(BINDATA_NAME)
	


