GO := $(shell which go)

SRC := $(wildcard **/*.go)
BIN := otpgen
DBFILE := users.db

$(BIN): $(SRC)
	$(GO) build -o $(BIN)

.PHONY: db
db:
	sqlite3 $(DBFILE) "CREATE TABLE IF NOT EXISTS users(username TEXT, seed TEXT, password TEXT, salt INTEGER);"

.PHONY: clean
clean:
	$(RM) $(BIN)


