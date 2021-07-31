GO := $(shell which go)
GOBIN := $(shell go env GOBIN)

SRC := $(wildcard **/*.go)
BIN := as
BINDIR := bin
DBFILE := clients.db
DEBUG_FLAGS := -gcflags="all=-N -l"

$(BIN): $(SRC)
	$(GO) build -o $(BINDIR)/$(BIN)

.PHONY: dev
dev: compile-daemon
	$(GOBIN)/CompileDaemon -directory=$(PWD) -exclude-dir=.git -build='make' -command='./bin/as' -color=true -log-prefix=false

.PHONY: db
db:
	sqlite3 $(DBFILE) "\
		CREATE TABLE IF NOT EXISTS clients(\
			id TEXT UNIQUE NOT NULL DEFAULT (lower(hex(randomblob(16)))),\
			secret_key TEXT NOT NULL DEFAULT (lower(hex(randomblob(8))))\
		);\
	"

.PHONY: get-clients
get-clients:
	@sqlite3 -header -column $(DBFILE) "SELECT * FROM clients"

.PHONY: add-clients
add-client:
	@sqlite3 -header -column $(DBFILE) "INSERT INTO clients DEFAULT VALUES "

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
