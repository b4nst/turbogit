package context

import (
	"fmt"

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

	return &Context{Username: viperOrGit("user", "name", r), Email: viperOrGit("user", "email", r), Repo: r}, nil
}

func viperOrGit(section string, option string, r *git.Repository) string {
	vk := fmt.Sprintf("%s.%s", section, option)
	value := viper.GetString(vk)

	if value == "" {
		cfg, err := r.Config()
		if err != nil {
			value = ""
		} else {
			value = cfg.Raw.Section(section).Option(option) // Switch to cfg.Merged when https://github.com/go-git/go-git/pull/20 is released
			viper.Set(vk, value)
		}
	}

	return value
}
