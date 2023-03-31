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

func TestStagedDiff(t *testing.T) {
	r := test.TestRepo(t)
	defer test.CleanupRepo(t, r)
	test.InitRepoConf(t, r)

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

	diff, err := StagedDiff(r)
	assert.NoError(t, err)
	deltas, err := diff.NumDeltas()
	assert.NoError(t, err)
	assert.Equal(t, 1, deltas)
}

func TestCurrentPatch(t *testing.T) {
	r := test.TestRepo(t)
	defer test.CleanupRepo(t, r)
	test.InitRepoConf(t, r)

	test.StageNewFile(t, r)
	c, err := Commit(r, "feat: initial commit")
	require.NoError(t, err)
	tree, err := c.Tree()
	require.NoError(t, err)

	test.StageNewFile(t, r)
	c1, err := Commit(r, "feat: second commit")
	require.NoError(t, err)
	tree1, err := c1.Tree()
	require.NoError(t, err)

	diff, err := r.DiffTreeToTree(tree, tree1, nil)
	require.NoError(t, err)

	s, err := PatchFromDiff(diff)
	fmt.Println(s)
	assert.NoError(t, err)
	assert.Contains(t, s, "new file mode 100644")
}
