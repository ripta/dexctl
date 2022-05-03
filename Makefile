BIN := dexctl

PREFIX ?= /usr
LIB_DIR = $(DESTDIR)$(PREFIX)/lib
BIN_DIR = $(DESTDIR)$(PREFIX)/bin
SHARE_DIR = $(DESTDIR)$(PREFIX)/share

export CGO_CPPFLAGS := ${CPPFLAGS}
export CGO_CFLAGS := ${CFLAGS}
export CGO_CXXFLAGS := ${CXXFLAGS}
export CGO_LDFLAGS := ${LDFLAGS}
export GOFLAGS := -buildmode=pie -trimpath -mod=readonly -modcacherw

.PHONY: local
local: vendor build

.PHONY: run
run: local
	go run main.go

.PHONY: build
build: main.go
	go build -o bin/$(BIN) main.go

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: clean
clean:
	rm -f "$(BIN)"
	rm -rf dist
	rm -rf vendor

.PHONY: fmt
fmt: ## Verifies all files have been `gofmt`ed.
	@echo "+ $@"
	@gofmt -s -l .

.PHONY: install
install:
	install -Dm755 -t "$(BIN_DIR)/" $(BIN)
