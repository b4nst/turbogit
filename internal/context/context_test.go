package context

import (
	"testing"

	"github.com/go-git/go-git/v5/plumbing/format/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestViperOrGit(t *testing.T) {
	r, teardown, err := setUp()
	defer teardown()
	if err != nil {
		t.Fatal(err)
	}

	// With git
	cfg, err := r.Config()
	if err != nil {
		t.Fatal(err)
	}
	cfg.Raw.AddOption("user", config.NoSubsection, "name", "bob")
	r.Storer.SetConfig(cfg)
	assert.Equal(t, "bob", viperOrGit("user", "name", r))

	// With viper
	viper.Set("user.name", "alice")
	assert.Equal(t, "alice", viperOrGit("user", "name", r))
}
