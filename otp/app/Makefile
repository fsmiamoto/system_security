GO := $(shell which go)

SRC := $(wildcard **/*.go)
BIN := app
BINDIR := bin
DBFILE := users.db
DEBUG_FLAGS := -gcflags="all=-N -l"

$(BIN): $(SRC)
	$(GO) build -o $(BINDIR)/$(BIN)

.PHONY: db
db:
	sqlite3 $(DBFILE) "\
		CREATE TABLE IF NOT EXISTS users(username TEXT UNIQUE NOT NULL, seed TEXT NOT NULL, salt TEXT NOT NULL);\
		CREATE TABLE IF NOT EXISTS invalid_otps(otp TEXT UNIQUE NOT NULL);\
	"

.PHONY: get-users
get-users:
	@sqlite3 -header -column $(DBFILE) "SELECT * FROM users"

.PHONY: get-otps
get-otps:
	@sqlite3 -header -column $(DBFILE) "SELECT * FROM invalid_otps"

.PHONY: clean
clean:
	$(RM) $(BIN)

.PHONY: clean-db
clean-db:
	$(RM) $(DBFILE)


