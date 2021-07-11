package user

import (
	"github.com/spf13/cobra"
)

var User = &cobra.Command{
	Use:   "user [add|rm]",
	Short: "User commands",
}

func init() {
	User.AddCommand(Add)
	User.AddCommand(Remove)
}
