package cmd

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/b4nst/turbogit/internal/format"
	"github.com/libgit2/git2go/v30"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBranchCreate(t *testing.T) {
	// Init git repository in tmp dir
	dir, err := ioutil.TempDir("", "turbogit-test-commit")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	require.NoError(t, os.Chdir(dir))
	r, err := git.InitRepository(dir, false)
	require.NoError(t, err)

	bco := &BranchCmdOption{
		format.FeatureBranch,
		"feat/foo",
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
	sig := &git.Signature{"alice@ecorp.com", "alice", time.Now()}
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
	assert.Equal(t, "refs/heads/"+bco.Name, h.Name())
}

func TestParseBranchCmd(t *testing.T) {
	dir, err := ioutil.TempDir("", "turbogit-test-commit")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	require.NoError(t, os.Chdir(dir))

	r, err := git.InitRepository(dir, false)
	require.NoError(t, err)
	config, err := r.Config()
	require.NoError(t, err)
	require.NoError(t, config.SetString("user.name", "alice"))

	cmd := &cobra.Command{}

	// Incorrect type
	_, err = parseBranchCmd(cmd, []string{"nope"})
	assert.Error(t, err, "Nope")
	// Not user no description
	_, err = parseBranchCmd(cmd, []string{"feat"})
	assert.Error(t, err, "feat branches need a description.")
	// User branch
	bco, err := parseBranchCmd(cmd, []string{"user"})
	assert.NoError(t, err)
	expected := BranchCmdOption{
		BType: format.UserBranch,
		Name:  "user/alice",
		Repo:  r,
	}
	assert.Equal(t, expected, *bco)
	// Classic branch
	bco, err = parseBranchCmd(cmd, []string{"feat", "Foo", "bAr"})
	assert.NoError(t, err)
	expected = BranchCmdOption{
		BType: format.FeatureBranch,
		Name:  "feat/foo-bar",
		Repo:  r,
	}
	assert.Equal(t, expected, *bco)
}
