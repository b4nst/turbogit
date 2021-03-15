package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	tugit "github.com/b4nst/turbogit/internal/git"
	"github.com/b4nst/turbogit/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunCheck(t *testing.T) {
	r := test.TestRepo(t)
	defer test.CleanupRepo(t, r)
	test.InitRepoConf(t, r)

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
