package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	tugit "github.com/b4nst/turbogit/internal/git"
	"github.com/b4nst/turbogit/internal/test"
	git "github.com/libgit2/git2go/v30"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunCheck(t *testing.T) {
	dir, err := ioutil.TempDir("", "turbogit-test-check")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	require.NoError(t, os.Chdir(dir))

	r, err := git.InitRepository(dir, false)
	require.NoError(t, err)
	config, err := r.Config()
	require.NoError(t, err)
	require.NoError(t, config.SetString("user.name", "alice"))
	require.NoError(t, config.SetString("user.email", "alice@ecorp.com"))

	c1, err := tugit.Commit(r, "bad commit 1")
	require.NoError(t, err)
	sid1, err := c1.ShortId()
	require.NoError(t, err)
	_, err = tugit.Commit(r, "feat: ok commit")
	assert.NoError(t, err)
	c3, err := tugit.Commit(r, "bad commit 2")
	assert.NoError(t, err)
	sid3, err := c3.ShortId()
	require.NoError(t, err)

	stderr, reset := test.CaptureStd(t, os.Stderr)
	err = runCheck(&CheckCmdOption{All: false, From: "HEAD", Repo: r})
	reset()
	assert.EqualError(t, err, "This commits are not compliant")
	stde, err := ioutil.ReadFile(stderr.Name())
	assert.Equal(t, fmt.Sprintf("%s %s\n%s %s\n", sid3, "bad commit 2", sid1, "bad commit 1"), string(stde))
}
