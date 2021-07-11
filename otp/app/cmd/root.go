package cmd

import (
	"fmt"
	"os"

	"github.com/fsmiamoto/system_security/otp/app/cmd/login"
	"github.com/fsmiamoto/system_security/otp/app/cmd/user"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "app",
	Short: "One-Time Password App",
}

func init() {
	root.AddCommand(user.User)
	root.AddCommand(login.Login)
}

func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
