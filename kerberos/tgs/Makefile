GO := $(shell which go)
GOBIN := $(shell go env GOBIN)
LDFLAGS = -X main.secretKey=$(SECRET_KEY) -X main.initVector=$(INIT_VECTOR)

SRC := $(wildcard **/*.go)
BIN := tgs
BINDIR := bin
DBFILE := services.db
DEBUG_FLAGS := -gcflags="all=-N -l"

# Ideally should be envvars but this is just for learning :)
SECRET_KEY = 5a6d29b8
INIT_VECTOR = 548ee275

$(BIN): $(SRC)
	$(GO) build -ldflags '$(LDFLAGS)' -o $(BINDIR)/$(BIN)

.PHONY: dev
dev: compile-daemon
	$(GOBIN)/CompileDaemon -directory=$(PWD) -exclude-dir=.git -build='make' -command='./bin/$(BIN)' -color=true -log-prefix=false

.PHONY: db
db:
	sqlite3 $(DBFILE) "\
		CREATE TABLE IF NOT EXISTS services(\
			id TEXT UNIQUE NOT NULL DEFAULT (lower(hex(randomblob(16)))),\
			secret_key TEXT NOT NULL DEFAULT (substring(lower(hex(randomblob(32))),8,8)),\
			iv TEXT NOT NULL DEFAULT (substring(lower(hex(randomblob(32))),8,8))\
		);\
	"

.PHONY: get-services
get-services:
	@sqlite3 -header -column $(DBFILE) "SELECT * FROM services"

.PHONY: add-services
add-service:
	@sqlite3 -header -column $(DBFILE) "INSERT INTO services DEFAULT VALUES "

.PHONY: clean
clean:
	$(RM) $(BINDIR)/$(BIN)

.PHONY: clean-db
clean-db:
	$(RM) $(DBFILE)

.PHONY: compile-daemon
compile-daemon:
	@if command -v "$(GOBIN)/CompileDaemon" > /dev/null; then\
		echo "CompileDaemon already present, skipping install...";\
	else\
		$(GO) install github.com/githubnemo/CompileDaemon@latest;\
	fi
