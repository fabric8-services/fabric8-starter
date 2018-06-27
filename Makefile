PROJECT_NAME=fabric8-starter
PACKAGE_NAME:=github.com/fabric8-services/$(PROJECT_NAME)

# -----------------------------------------------------------------------
# bootstrap
# catch 22: can't include from vendor deps before running the `deps` goal, 
# but `deps` is defined in `Makefile.common`!
# ------------------------------------------------------------------------
CUR_DIR=$(shell pwd)
MAKEFILE_DIR=$(CUR_DIR)/tmp/makefile
include $(MAKEFILE_DIR)/Makefile.common

makefiles: $(MAKEFILE_DIR)/Makefile.common $(MAKEFILE_DIR)/Makefile.lnx $(MAKEFILE_DIR)/Makefile.win
$(MAKEFILE_DIR):
	@mkdir -p $(MAKEFILE_DIR)

$(MAKEFILE_DIR)/Makefile.common: $(MAKEFILE_DIR)
	@curl -o $(MAKEFILE_DIR)/Makefile.common https://raw.githubusercontent.com/xcoulon/fabric8-common/base_deps/makefile/Makefile.common

$(MAKEFILE_DIR)/Makefile.lnx: $(MAKEFILE_DIR)
	@curl -o $(MAKEFILE_DIR)/Makefile.lnx https://raw.githubusercontent.com/xcoulon/fabric8-common/base_deps/makefile/Makefile.lnx

$(MAKEFILE_DIR)/Makefile.win: $(MAKEFILE_DIR)
	@curl -o $(MAKEFILE_DIR)/Makefile.win https://raw.githubusercontent.com/xcoulon/fabric8-common/base_deps/makefile/Makefile.win

# -------------------------------------------------------------------
# run in dev mode
# -------------------------------------------------------------------
.PHONY: dev
dev: prebuild-check deps generate $(FRESH_BIN) ## run the server locally
	F8_DEVELOPER_MODE_ENABLED=true $(FRESH_BIN)

# -------------------------------------------------------------------
# build the binary executable (to ship in prod)
# -------------------------------------------------------------------
LDFLAGS=-ldflags "-X ${PACKAGE_NAME}/controller.Commit=${COMMIT} -X ${PACKAGE_NAME}/controller.BuildTime=${BUILD_TIME}"

.PHONY: build
build: makefiles prebuild-check deps generate ## Build the server
ifeq ($(OS),Windows_NT)
	go build -v $(LDFLAGS) -o "$(shell cygpath --windows '$(BINARY_SERVER_BIN)')"
else
	go build -v $(LDFLAGS) -o $(BINARY_SERVER_BIN)
endif

.PHONY: generate
generate: prebuild-check $(DESIGNS) $(GOAGEN_BIN) $(VENDOR_DIR) ## Generate GOA sources. Only necessary after clean of if changed `design` folder.
	$(GOAGEN_BIN) app -d ${PACKAGE_NAME}/${DESIGN_DIR}
	$(GOAGEN_BIN) controller -d ${PACKAGE_NAME}/${DESIGN_DIR} -o controller/ --pkg controller --app-pkg ${PACKAGE_NAME}/app
	$(GOAGEN_BIN) gen -d ${PACKAGE_NAME}/${DESIGN_DIR} --pkg-path=github.com/fabric8-services/fabric8-common/goasupport/status --out app
	$(GOAGEN_BIN) gen -d ${PACKAGE_NAME}/${DESIGN_DIR} --pkg-path=github.com/fabric8-services/fabric8-common/goasupport/jsonapi_errors_helpers --out app
	$(GOAGEN_BIN) swagger -d ${PACKAGE_NAME}/${DESIGN_DIR}
	