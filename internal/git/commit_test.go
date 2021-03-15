package git

import (
	"testing"
	"time"

	"github.com/b4nst/turbogit/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommit(t *testing.T) {
	r := test.TestRepo(t)
	defer test.CleanupRepo(t, r)
	test.InitRepoConf(t, r)

	commit, err := Commit(r, "commit message")
	assert.NoError(t, err)
	assert.Equal(t, "commit message", commit.Message())
	assert.Equal(t, test.GIT_USERNAME, commit.Author().Name)
	assert.Equal(t, test.GIT_EMAIL, commit.Author().Email)
	assert.WithinDuration(t, time.Now(), commit.Author().When, 5*time.Second)
	head, err := r.Head()
	require.NoError(t, err)
	headCommit, err := r.LookupCommit(head.Target())
	require.NoError(t, err)
	assert.Equal(t, headCommit.Id(), commit.Id())
}

func TestAmend(t *testing.T) {
	r := test.TestRepo(t)
	test.InitRepoConf(t, r)
	commit, err := Commit(r, "foo")
	require.NoError(t, err)

	amendc, err := Amend(commit, "bar")
	assert.NoError(t, err)
	head, err := r.Head()
	require.NoError(t, err)
	headc, err := r.LookupCommit(head.Target())
	require.NoError(t, err)
	assert.Equal(t, "bar", headc.Message())
	assert.Equal(t, headc.Id(), amendc.Id())
}
