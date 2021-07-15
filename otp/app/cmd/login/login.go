package login

import (
	"errors"

	"github.com/fsmiamoto/system_security/otp/app/otp"
	"github.com/fsmiamoto/system_security/otp/app/repository"
	"github.com/spf13/cobra"
)

var ErrAccessDenied = errors.New("access denied")

var Login = &cobra.Command{
	Use:          "login [username] [otp]",
	Short:        "Login into the app",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return cmd.Help()
		}
		return login(args[0], args[1])
	},
}

func login(username, password string) error {
	user, err := repository.Get(username)
	if err != nil {
		return err
	}

	if repository.Exists(password) {
		return ErrAccessDenied
	}

	// TODO: Store only the hash of both seed and salt
	otps := otp.NewList(5, user.Seed, user.Salt)

	index, allow := 0, false
	for i := range otps {
		if otps[i] == password {
			index, allow = i, true
		}
	}

	if !allow {
		return ErrAccessDenied
	}

	repository.Invalidate(otps[index:])

	return nil
}
