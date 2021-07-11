package cmd

import (
	"fmt"
	"os"

	"github.com/fsmiamoto/system_security/otp/otpgen/cmd/gen"
	"github.com/fsmiamoto/system_security/otp/otpgen/cmd/user"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "otpgen",
	Short: "One-Time Password Generator",
}

func init() {
	root.AddCommand(user.User)
	root.AddCommand(gen.Gen)
}

func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
