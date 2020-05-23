package context

import (
	"testing"

	"github.com/go-git/go-git/plumbing/format/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestUsername(t *testing.T) {
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
	assert.Equal(t, "bob", username(r))

	// With viper
	viper.Set("user.name", "alice")
	assert.Equal(t, "alice", username(r))
}

func TestEmail(t *testing.T) {
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
	cfg.Raw.AddOption("user", config.NoSubsection, "email", "bob@company.com")
	r.Storer.SetConfig(cfg)
	assert.Equal(t, "bob@company.com", email(r))

	// With viper
	viper.Set("user.email", "alice@company.com")
	assert.Equal(t, "alice@company.com", email(r))
}
