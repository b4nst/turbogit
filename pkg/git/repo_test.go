package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/b4nst/turbogit/pkg/test"
	git "github.com/libgit2/git2go/v33"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetrepo(t *testing.T) {
	dir, err := ioutil.TempDir("", "turbogit-test-getrepo")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	require.NoError(t, os.Chdir(dir))

	// Working dir is not a repo
	_, err = Getrepo()
	assert.Error(t, err)
	if err != nil {
		assert.Contains(t, err.Error(), "could not find repository from")
	}

	// Working dir is a repo
	r, err := git.InitRepository(dir, false)
	require.NoError(t, err)
	repo, err := Getrepo()
	assert.NoError(t, err)
	assert.Equal(t, r, repo)
}

func TestCurrentPatch(t *testing.T) {
	r := test.TestRepo(t)
	defer test.CleanupRepo(t, r)

	f := test.NewFile(t, r)
	test.StageFile(t, f, r)
	_, err := Commit(r, "feat: initial commit")
	require.NoError(t, err)
	// Staged stuff
	fmt.Fprintln(f, "Staged")
	test.StageFile(t, f, r)
	// Not staged stuff
	fmt.Fprintln(f, "Not_staged")
	test.NewFile(t, r)

	s, err := CurrentPatch(r)
	assert.NoError(t, err)
	assert.Contains(t, s, "+Staged")
	assert.NotContains(t, s, "Not_staged")
}
