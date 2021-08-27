GO := $(shell which go)
GOBIN := $(shell go env GOBIN)

SRC := $(wildcard **/*.go)
BIN := http_proxy
BINDIR := bin
DEBUG_FLAGS := -gcflags="all=-N -l"

$(BIN): $(SRC)
	$(GO) build -ldflags '$(LDFLAGS)' -o $(BINDIR)/$(BIN)

.PHONY: dev
dev: compile-daemon
	$(GOBIN)/CompileDaemon -directory=$(PWD) -exclude-dir=.git -build='make' -command='./bin/$(BIN)' -color=true -log-prefix=false

.PHONY: clean
clean:
	$(RM) $(BINDIR)/$(BIN)

.PHONY: compile-daemon
compile-daemon:
	@if command -v "$(GOBIN)/CompileDaemon" > /dev/null; then\
		echo "CompileDaemon already present, skipping install...";\
	else\
		$(GO) install github.com/githubnemo/CompileDaemon@latest;\
	fi

.PHONY: checksum
checksum:
	md5sum $(BINDIR)/$(BIN)
