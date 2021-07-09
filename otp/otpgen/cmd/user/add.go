package user

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Add = &cobra.Command{
	Use:   "add [user] [password]",
	Short: "Add a new user to the database",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Help()
			return
		}

		fmt.Println("User: " + args[0])
		fmt.Println("Password: " + args[1])
	},
}
