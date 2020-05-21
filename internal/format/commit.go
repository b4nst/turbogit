package format

import (
	"fmt"

	"github.com/b4nst/turbogit/internal/constants"
)

type CommitType int

const (
	BuildCommit CommitType = iota
	CiCommit
	ChoreCommit
	DocCommit
	FeatureCommit
	FixCommit
	NilCommit
	PerfCommit
	RefactorCommit
	StyleCommit
	TestCommit
)

func (b CommitType) String() string {
	return [...]string{"build", "ci", "chore", "docs", "feat", "fix", "", "perf", "refactor", "style", "test"}[b]
}

type CommitMessageOption struct {
	// Commit type (optional)
	Ctype CommitType
	// Commit scope (optional)
	Scope string
	// Commit subject (required)
	Description string
	// Commit body (optional)
	Body string
	// Commit footers (optional)
	Footers []string
	// Breaking change flag (optional)
	BreakingChanges bool
}

// Format commit message according to https://www.conventionalcommits.org/en/v1.0.0/
func CommitMessage(o *CommitMessageOption) string {
	msg := o.Ctype.String()
	// Add scope if any
	if o.Scope != "" {
		msg += fmt.Sprintf("(%s)", o.Scope)
	}
	// Mark breaking changes
	if o.BreakingChanges {
		msg += "!"
	}
	// Add description
	msg += fmt.Sprintf(": %s", o.Description)
	// Add body if any
	if o.Body != "" {
		msg += constants.LINE_BREAK + constants.LINE_BREAK + o.Body
	}
	// Add footers if any
	if len(o.Footers) > 0 {
		msg += constants.LINE_BREAK
		for _, f := range o.Footers {
			msg += constants.LINE_BREAK + f
		}
	}

	return msg
}
