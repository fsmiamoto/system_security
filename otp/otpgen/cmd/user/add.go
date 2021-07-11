package user

import (
	"github.com/fsmiamoto/system_security/otp/otpgen/repository"
	"github.com/spf13/cobra"
)

var Add = &cobra.Command{
	Use:   "add [username] [password]",
	Short: "Add a new user to the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return cmd.Help()
		}
		return add(args[0], args[1])
	},
}

var Remove = &cobra.Command{
	Use:   "rm [username]",
	Short: "Remove a user from the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Help()
		}
		return remove(args[0])
	},
}

func add(username, password string) error {
	return repository.Add(username, password)
}

func remove(username string) error {
	return repository.Remove(username)
}
