GO := $(shell which go)

SRC := $(wildcard **/*.go)
BIN := otpgen
BINDIR := bin
DBFILE := users.db

$(BIN): $(SRC)
	$(GO) build -o $(BINDIR)/$(BIN)

.PHONY: db
db:
	sqlite3 $(DBFILE) "\
		CREATE TABLE IF NOT EXISTS users(username TEXT, seed TEXT, password TEXT, salt TEXT);\
		CREATE TABLE IF NOT EXISTS otp(password TEXT, used BOOLEAN);\
	"

.PHONY: get-users
get-users:
	@sqlite3 -header -column $(DBFILE) "SELECT * FROM users"

.PHONY: clean
clean:
	$(RM) $(BIN)

.PHONY: clean-db
clean-db:
	$(RM) $(DBFILE)


