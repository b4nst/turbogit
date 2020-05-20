package format

import (
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
			expected: "features/my-branch",
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
			expected: "users/bob",
		},
		"User branch with description": {
			t:        UserBranch,
			d:        "alice branch",
			u:        "alice",
			expected: "users/alice/alice-branch",
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, BranchName(tc.t, tc.d, tc.u))
		})
	}
}
