package config

import (
	"fmt"

	git "github.com/libgit2/git2go/v30"
)

type AnyBoolDirective struct {
	Key    string
	Config *git.Config
}

func (d *AnyBoolDirective) Apply(autofix bool) error {
	if _, err := d.Config.LookupBool(d.Key); err == nil {
		return nil
	}
	if autofix {
		d.Config.SetBool(d.Key, false)
		return nil
	}
	return fmt.Errorf("%s is not a boolean", d.Key)
}

type BoolDirective struct {
	Key      string
	Expected bool
	Config   *git.Config
}

func (d *BoolDirective) Apply(autofix bool) error {
	if v, err := d.Config.LookupBool(d.Key); err == nil && v == d.Expected {
		return nil
	}
	if autofix {
		d.Config.SetBool(d.Key, d.Expected)
		return nil
	}
	return fmt.Errorf("%s is not a boolean equal to %v", d.Key, d.Expected)
}
