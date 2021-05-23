package config

import (
	"errors"

	git "github.com/libgit2/git2go/v30"
)

type Directive interface {
	Apply(bool) error
}

type DirectiveConf struct {
	Name        string
	Description string
	Key         string
	Value       interface{}
	Scope       string
	Only        string
}

func (df *DirectiveConf) Compile(r *git.Repository) (*Directive, error) {
	return nil, errors.New("Not implemented")
}
