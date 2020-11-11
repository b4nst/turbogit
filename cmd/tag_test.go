package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/b4nst/turbogit/internal/test"
	"github.com/blang/semver/v4"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLastTag(t *testing.T) {
	r, teardown, err := test.SetUpRepo()
	defer teardown()
	require.NoError(t, err)

	// No tag
	last, err := lastTag(r)
	require.NoError(t, err)
	assert.Nil(t, last)

	testNewVersion := func(v string) {
		head, err := r.Head()
		require.NoError(t, err)
		ref, err := r.CreateTag(v, head.Hash(), nil)
		require.NoError(t, err)
		err = r.Storer.SetReference(ref)
		require.NoError(t, err)
		last, err = lastTag(r)
		require.NoError(t, err)
		require.NotNil(t, last)
		assert.Equal(t, ref, last.ref)
		assert.Equal(t, v, last.version.String())
	}

	testNewVersion("1.0.0-alpha")
	testNewVersion("1.0.0-alpha.1")
	testNewVersion("1.0.0-beta")
	testNewVersion("1.0.0-beta.1")
	testNewVersion("1.0.1")
	testNewVersion("1.1.0")
	testNewVersion("2.0.0")
}

func TestNextVersion(t *testing.T) {
	tcs := map[string]struct {
		msgs     []string
		curr     semver.Version
		expected semver.Version
	}{
		"Same":  {msgs: []string{"malformatted", "refactor: some refactoring"}, curr: semver.MustParse("1.0.0"), expected: semver.MustParse("1.0.0")},
		"Patch": {msgs: []string{"fix: msg", "refactor: ref"}, curr: semver.MustParse("1.0.0"), expected: semver.MustParse("1.0.1")},
		"Minor": {msgs: []string{"fix: msg", "feat: msg"}, curr: semver.MustParse("1.0.0"), expected: semver.MustParse("1.1.0")},
		"Major": {msgs: []string{"fix: msg", "feat: msg", "build!: msg"}, curr: semver.MustParse("1.0.0"), expected: semver.MustParse("2.0.0")},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			actual, err := nextVersion(tc.curr, tc.msgs)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestCommitMsgsSince(t *testing.T) {
	r, teardown, err := test.SetUpRepo()
	defer teardown()
	require.NoError(t, err)
	wt, err := r.Worktree()
	require.NoError(t, err)

	test.StageNewFile(r)
	h, err := wt.Commit("tag commit", &git.CommitOptions{})
	require.NoError(t, err)
	ref, err := r.CreateTag("tag", h, nil)

	test.StageNewFile(r)
	_, err = wt.Commit("msg 1", &git.CommitOptions{})
	require.NoError(t, err)
	test.StageNewFile(r)
	_, err = wt.Commit("msg 2", &git.CommitOptions{})
	require.NoError(t, err)
	test.StageNewFile(r)
	_, err = wt.Commit("msg 3", &git.CommitOptions{})
	require.NoError(t, err)

	actual, err := commitMsgsSince(r, ref.Hash())
	assert.NoError(t, err)
	assert.Equal(t, []string{"msg 3", "msg 2", "msg 1"}, actual)
}

func TestTag(t *testing.T) {
	r, teardown, err := test.SetUpRepo()
	defer teardown()
	require.NoError(t, err)
	wt, err := r.Worktree()
	require.NoError(t, err)

	stdout := os.Stdout
	restore := func() {
		os.Stdout = stdout
	}
	defer restore()
	devnull, err := ioutil.TempFile("", "dev-null")
	require.NoError(t, err)
	defer os.RemoveAll(devnull.Name())
	os.Stdout = devnull

	cmd := &cobra.Command{}
	dr := cmd.Flags().BoolP("dry-run", "d", false, "Do not tag.")

	test.StageNewFile(r)
	h, err := wt.Commit("feat: first feat", &git.CommitOptions{})
	require.NoError(t, err)
	*dr = true
	assert.NoError(t, tag(cmd, []string{}))
	lt, err := test.LastTagFrom(r)
	require.NoError(t, err)
	assert.Nil(t, lt)

	_, err = r.CreateTag("tag", h, nil) // In case of not semver tag
	*dr = false
	assert.NoError(t, tag(cmd, []string{}))
	lt, err = test.LastTagFrom(r)
	require.NoError(t, err)
	assert.Equal(t, "v0.1.0", lt.Name().Short())

	test.StageNewFile(r)
	_, err = wt.Commit("fix!: breaking changes", &git.CommitOptions{})
	require.NoError(t, err)
	*dr = false
	assert.NoError(t, tag(cmd, []string{}))
	lt, err = test.LastTagFrom(r)
	require.NoError(t, err)
	assert.Equal(t, "v1.0.0", lt.Name().Short())
}
