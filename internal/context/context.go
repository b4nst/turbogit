package context

import (
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Context struct {
	// Git user name
	Username string

	// Git user email
	Email string

	// Repository
	Repo *git.Repository
}

func FromCommand(cmd *cobra.Command) (*Context, error) {
	r, err := currRepo()
	if err != nil {
		return nil, err
	}

	return &Context{Username: getUsername(r), Email: getEmail(r), Repo: r}, nil
}

func getUsername(r *git.Repository) string {
	username := viper.GetString("username")

	if username == "" {
		cfg, err := r.Config()
		if err != nil {
			username = ""
		} else {
			username = cfg.Raw.Section("user").Option("name") // Switch to cfg.Merged when https://github.com/go-git/go-git/pull/20 is released
		}
	}

	return username
}

func getEmail(r *git.Repository) string {
	email := viper.GetString("email")

	if email == "" {
		cfg, err := r.Config()
		if err != nil {
			email = ""
		} else {
			email = cfg.Raw.Section("user").Option("email") // Switch to cfg.Merged when https://github.com/go-git/go-git/pull/20 is released
		}
	}

	return email
}
