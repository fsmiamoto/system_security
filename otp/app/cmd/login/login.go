package login

import "github.com/spf13/cobra"

var Login = &cobra.Command{
	Use:   "login [username] [otp]",
	Short: "Login into the app",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return cmd.Help()
		}
		return login(args[0], args[1])
	},
}

func login(username, otp string) error {
	return nil
}
