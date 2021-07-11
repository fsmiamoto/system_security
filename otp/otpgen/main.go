package main

import (
	"math/rand"
	"time"

	"github.com/fsmiamoto/system_security/otp/otpgen/cmd"
)

func main() {
	rand.Seed(time.Now().Unix())
	cmd.Execute()
}
