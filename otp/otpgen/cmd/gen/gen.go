package gen

import (
	"database/sql"
	"errors"

	"github.com/fsmiamoto/system_security/otp/otpgen/hash"
	"github.com/fsmiamoto/system_security/otp/otpgen/repository"
	"github.com/spf13/cobra"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

var Gen = &cobra.Command{
	Use:   "gen [username] [password]",
	Short: "Generate a list of one-time passwords for the user",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return cmd.Help()
		}

		user, err := repository.Get(args[0])

		switch err {
		case sql.ErrNoRows:
			return errors.New("invalid credentials")
		case nil:
			break
		default:
			return err
		}

		hashedPasswordWithSalt := hash.Sha256(args[1] + user.Salt)

		if hashedPasswordWithSalt != user.Password {
			return errors.New("invalid credentials")
		}

		return nil
	},
}
