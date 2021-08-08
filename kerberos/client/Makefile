GO := $(shell which go)
GOBIN := $(shell go env GOBIN)
LDFLAGS = -X main.clientID=$(CLIENT_ID) -X main.secretKey=$(SECRET_KEY) -X main.initVector=$(INIT_VECTOR)

SRC := $(wildcard **/*.go)
BIN := client
BINDIR := bin
DEBUG_FLAGS := -gcflags="all=-N -l"

# Ideally should be envvars but this is just for learning :)
CLIENT_ID = 34e15fc03b81928f9f4386adc6e45fa9
SECRET_KEY = 2963252b 
INIT_VECTOR = cafc7a3b

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
