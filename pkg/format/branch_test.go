package format

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTugBranchString(t *testing.T) {
	tcs := map[string]struct {
		tb       TugBranch
		expected string
	}{
		"Without prefix": {
			tb:       TugBranch{Type: "feat", Description: "A foo feature."},
			expected: "feat/a-foo-feature",
		},
		"With prefix": {
			tb:       TugBranch{Type: "user", Prefix: "alice", Description: "Alice branch"},
			expected: "user/alice/alice-branch",
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.tb.String())
		})
	}
}

func TestParseBranch(t *testing.T) {
	tcs := map[string]struct {
		str      string
		err      error
		expected TugBranch
	}{
		"Without prefix": {
			str:      "feat/a-foo-feature",
			err:      nil,
			expected: TugBranch{Type: "feat", Description: "A foo feature"},
		},
		"With prefix": {
			str:      "user/alice/alice-branch",
			err:      nil,
			expected: TugBranch{Type: "user", Prefix: "alice", Description: "Alice branch"},
		},
		"Error branch": {
			str:      "BADBEEF",
			err:      errors.New("Bad branch format"),
			expected: TugBranch{},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			b, err := ParseBranch(tc.str)
			assert.Equal(t, tc.expected, b)
			if tc.err == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.err.Error())
			}
		})
	}
}

func TestTugBranchWithType(t *testing.T) {
	tcs := map[string]struct {
		t        string
		rw       map[string]string
		expected TugBranch
	}{
		"Empty rewrite": {
			t:        "type",
			rw:       map[string]string{},
			expected: TugBranch{Type: "type"},
		},
		"With rewrite": {
			t:        "type",
			rw:       map[string]string{"type": "foo"},
			expected: TugBranch{Type: "foo"},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			tb := TugBranch{}.WithType(tc.t, tc.rw)
			assert.Equal(t, tc.expected, tb)
		})
	}
}

func TestSanitizeBranch(t *testing.T) {
	dirty := "/A dirty branch/`should` ?be *~^:\\ [cleaned]../"

	assert.Equal(t, "A-dirty-branch/should-be-cleaned", sanitizeBranch(dirty))
}
