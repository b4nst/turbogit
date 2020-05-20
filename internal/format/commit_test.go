package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommitMessage(t *testing.T) {
	tcs := map[string]struct {
		o        *CommitMessageOption
		expected string
	}{
		"Simple subject feature": {
			o:        &CommitMessageOption{Ctype: FeatureCommit, Description: "commit description"},
			expected: "feat: commit description",
		},
		"Subject + scope perf": {
			o:        &CommitMessageOption{Ctype: PerfCommit, Description: "message", Scope: "scope"},
			expected: "perf(scope): message",
		},
		"Breaking change refactor": {
			o:        &CommitMessageOption{Ctype: RefactorCommit, Description: "message", BreakingChanges: true},
			expected: "refactor!: message",
		},
		"Full stuff": {
			o:        &CommitMessageOption{Ctype: FeatureCommit, Scope: "scope", Description: "message", BreakingChanges: true, Body: "The message body", Footers: []string{"First foot", "Second foot"}},
			expected: "feat(scope)!: message\n\nThe message body\n\nFirst foot\nSecond foot",
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, CommitMessage(tc.o))
		})
	}
}
