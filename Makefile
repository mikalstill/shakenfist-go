# Shaken Fist Golang Client makefile
#
# For linter installation see
#     https://golangci-lint.run/usage/install/#ci-installation
#

GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
TEST?=./...
TEST_COUNT?=1

default: build


build: fmtcheck
	go build .

test: fmtcheck
	go test $(TEST) $(TESTARGS) -v -timeout=120s -parallel=4

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./$(PKG_NAME)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/check_files_gofmt.sh'"

lint:
	@golangci-lint run ./$(PKG_NAME)/...

install-tools:
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint


.PHONY: build lint test fmt fmtcheck lint install-tools
