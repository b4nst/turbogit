package cmd

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/b4nst/turbogit/internal/format"
	git "github.com/libgit2/git2go/v30"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBranchCreate(t *testing.T) {
	// Init git repository in tmp dir
	dir, err := ioutil.TempDir("", "turbogit-test-branch")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	require.NoError(t, os.Chdir(dir))
	r, err := git.InitRepository(dir, false)
	require.NoError(t, err)

	bco := &BranchCmdOption{
		format.TugBranch{Type: "feat", Description: "foo"},
		r,
	}
	// Sanity test
	err = runBranch(bco)
	assert.Error(t, err, "No commit to create branch from, please create the initial commit")
	// Actually create branch
	config, err := r.Config()
	require.NoError(t, err)
	require.NoError(t, config.SetString("user.name", "alice"))
	require.NoError(t, config.SetString("user.email", "alice@ecorp.com"))
	sig := &git.Signature{
		Email: "alice@ecorp.com",
		Name:  "alice",
		When:  time.Now(),
	}
	idx, err := r.Index()
	require.NoError(t, err)
	treeId, err := idx.WriteTree()
	require.NoError(t, err)
	tree, err := r.LookupTree(treeId)
	require.NoError(t, err)
	_, err = r.CreateCommit("HEAD", sig, sig, "Initial commit", tree)

	err = runBranch(bco)
	assert.NoError(t, err)
	h, err := r.Head()
	require.NoError(t, err)
	assert.Equal(t, "refs/heads/feat/foo", h.Name())
}

func TestParseBranchCmd(t *testing.T) {
	dir, err := ioutil.TempDir("", "turbogit-test-branch")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	require.NoError(t, os.Chdir(dir))

	r, err := git.InitRepository(dir, false)
	require.NoError(t, err)
	config, err := r.Config()
	require.NoError(t, err)
	require.NoError(t, config.SetString("user.name", "alice"))

	cmd := &cobra.Command{}

	// User branch
	bco, err := parseBranchCmd(cmd, []string{"user", "my", "branch"})
	assert.NoError(t, err)
	expected := BranchCmdOption{
		NewBranch: format.TugBranch{Type: "user", Prefix: "alice", Description: "my branch"},
		Repo:      r,
	}
	assert.Equal(t, expected, *bco)
	// Classic branch
	bco, err = parseBranchCmd(cmd, []string{"feat", "foo", "bar"})
	assert.NoError(t, err)
	expected = BranchCmdOption{
		NewBranch: format.TugBranch{Type: "feat", Description: "foo bar"},
		Repo:      r,
	}
	assert.Equal(t, expected, *bco)
}
