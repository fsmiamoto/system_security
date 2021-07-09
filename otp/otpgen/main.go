package main

import (
	"math/rand"
	"time"

	"github.com/fsmiamoto/system_security/otp/otpgen/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	rand.Seed(time.Now().Unix())
	cmd.Execute()
}
