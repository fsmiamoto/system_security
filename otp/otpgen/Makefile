GO := $(shell which go)

SRC := $(wildcard **/*.go)
BIN := otpgen

$(BIN): $(SRC)
	$(GO) build -o $(BIN)

.PHONY: clean
clean:
	$(RM) $(BIN)


