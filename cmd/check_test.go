package cmd

import (
	"fmt"
	"testing"

	tugit "github.com/b4nst/turbogit/pkg/git"
	"github.com/b4nst/turbogit/pkg/test"
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

	err = check(&checkOpt{All: false, From: "HEAD", Repo: r})
	assert.EqualError(t, err, fmt.Sprintf("2 errors occurred:\n\t* %s ('bad commit 2') is not compliant\n\t* %s ('bad commit 1') is not compliant\n\n", sid3, sid1))
}
