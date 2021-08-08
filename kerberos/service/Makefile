GO := $(shell which go)
GOBIN := $(shell go env GOBIN)
LDFLAGS = -X main.secretKey=$(SECRET_KEY) -X main.initVector=$(INIT_VECTOR)

SRC := $(wildcard **/*.go)
BIN := service
BINDIR := bin
DBFILE := services.db
DEBUG_FLAGS := -gcflags="all=-N -l"

# Ideally should be envvars but this is just for learning :)
SECRET_KEY = aadf4fce
INIT_VECTOR = 38c08650

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
