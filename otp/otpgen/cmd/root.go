package cmd

import (
	"fmt"
	"os"

	"github.com/fsmiamoto/system_security/otp/otpgen/cmd/user"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "otpgen",
	Short: "One-Time Password Generator",
}

func init() {
	rootCmd.AddCommand(user.User)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
