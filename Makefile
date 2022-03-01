# Shaken Fist Golang Client makefile
#
# For linter installation see
#     https://golangci-lint.run/usage/install/#ci-installation
#

EXAMPLE_DIRS=$(wildcard examples/*)
EXAMPLES=$(foreach dir, $(EXAMPLE_DIRS), $(addsuffix /$(notdir $(dir)), $(dir)))
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
TEST?=$$(go list ./... | grep -v /examples/)

default: build


build: fmtcheck
	go build .

build-examples: $(EXAMPLES)

$(EXAMPLES):
	go build -o $@ .

clean:
	rm -f $(EXAMPLES)

test: fmtcheck
	go test $(TEST) $(TESTARGS) -v -timeout=120s -parallel=4

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w .

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/check_files_gofmt.sh'"

lint:
	golangci-lint run ./...

install-tools:
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.44.2

pre-commit: fmtcheck lint clean build-examples

.PHONY: build lint test fmt fmtcheck lint install-tools pre-commit
