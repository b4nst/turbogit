package cmd

import (
	"testing"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLastTag(t *testing.T) {
	r, teardown, err := setUpRepo()
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
