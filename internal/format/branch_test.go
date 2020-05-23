package format

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeBranch(t *testing.T) {
	dirty := "/A dirty branch/should ?be *~^:\\ cleaned../"

	assert.Equal(t, "a-dirty-branch/should-be-cleaned", sanitizeBranch(dirty))
}

func TestBranchName(t *testing.T) {
	tcs := map[string]struct {
		t        BranchType
		d        string
		u        string
		expected string
	}{
		"Feature branch": {
			t:        FeatureBranch,
			d:        "my branch",
			u:        "bob",
			expected: "feat/my-branch",
		},
		"Fix branch": {
			t:        FixBranch,
			d:        "some patch",
			u:        "alice",
			expected: "fix/some-patch",
		},
		"User branch": {
			t:        UserBranch,
			d:        "",
			u:        "bob",
			expected: "user/bob",
		},
		"User branch with description": {
			t:        UserBranch,
			d:        "alice branch",
			u:        "alice",
			expected: "user/alice/alice-branch",
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, BranchName(tc.t, tc.d, tc.u))
		})
	}
}

func TestBranchTypeFrom(t *testing.T) {
	tcs := map[string]struct {
		str      string
		err      error
		expected BranchType
	}{
		"Feature branch": {
			str:      "feat",
			err:      nil,
			expected: FeatureBranch,
		},
		"Fix branch": {
			str:      "fix",
			err:      nil,
			expected: FixBranch,
		},
		"User branch": {
			str:      "user",
			err:      nil,
			expected: UserBranch,
		},
		"Error branch": {
			str:      "nope",
			err:      errors.New("error"),
			expected: -1,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			b, err := BranchTypeFrom(tc.str)
			assert.Equal(t, tc.expected, b)
			if tc.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
