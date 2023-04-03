	"fmt"
	"github.com/b4nst/turbogit/pkg/test"

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