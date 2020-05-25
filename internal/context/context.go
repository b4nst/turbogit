package context

import (
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

type Context struct {
	// Repository
	Repo *git.Repository
}

func FromCommand(cmd *cobra.Command) (*Context, error) {
	r, err := currRepo()
	if err != nil {
		return nil, err
	}

	return &Context{Repo: r}, nil
}
