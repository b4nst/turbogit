package format

import (
	"errors"
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

func TestCMOCheck(t *testing.T) {
	tcs := map[string]struct {
		cmo *CommitMessageOption
		err error
	}{
		"No type": {
			cmo: &CommitMessageOption{},
			err: errors.New("A commit type is required"),
		},
		"No description": {
			cmo: &CommitMessageOption{Ctype: FeatureCommit},
			err: errors.New("A commit description is required"),
		},
		"Ok": {
			cmo: &CommitMessageOption{Ctype: FeatureCommit, Description: "foo"},
			err: nil,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.err, tc.cmo.Check())
		})
	}
}

func TestCMOOverwrite(t *testing.T) {
	tcs := map[string]struct {
		src      *CommitMessageOption
		override *CommitMessageOption
		expected *CommitMessageOption
	}{
		"Override type": {
			src:      &CommitMessageOption{Ctype: FeatureCommit, Description: "foo"},
			override: &CommitMessageOption{Ctype: FixCommit},
			expected: &CommitMessageOption{Ctype: FixCommit, Description: "foo"},
		},
		"Override everything": {
			src:      &CommitMessageOption{Ctype: FeatureCommit, Description: "foo", Scope: "foo", Body: "foo", Footers: []string{"foo", "foo"}, BreakingChanges: true},
			override: &CommitMessageOption{Ctype: FixCommit, Description: "bar", Scope: "bar", Body: "bar", Footers: []string{"bar", "bar"}, BreakingChanges: false},
			expected: &CommitMessageOption{Ctype: FixCommit, Description: "bar", Scope: "bar", Body: "bar", Footers: []string{"bar", "bar"}, BreakingChanges: true},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			err := tc.src.Overwrite(tc.override)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, tc.src)
		})
	}
}

func TestFindCommitType(t *testing.T) {
	tcs := map[string]struct {
		str      string
		expected CommitType
	}{
		"Nil": {"fail", NilCommit},

		"B":      {"b", BuildCommit},
		"Build":  {"bUild", BuildCommit},
		"Builds": {"builds", BuildCommit},

		"Ci": {"ci", CiCommit},

		"Ch":     {"ch", ChoreCommit},
		"Chore":  {"chore", ChoreCommit},
		"Chores": {"chOreS", ChoreCommit},

		"D":    {"d", DocCommit},
		"Doc":  {"Doc", DocCommit},
		"Docs": {"docs", DocCommit},

		"Fe":       {"fe", FeatureCommit},
		"Feat":     {"feAt", FeatureCommit},
		"Feats":    {"feats", FeatureCommit},
		"Feature":  {"feature", FeatureCommit},
		"Features": {"features", FeatureCommit},

		"Fi":    {"fi", FixCommit},
		"Fix":   {"Fix", FixCommit},
		"Fixes": {"fixEs", FixCommit},

		"P":            {"p", PerfCommit},
		"Perf":         {"perf", PerfCommit},
		"Perfs":        {"pErFs", PerfCommit},
		"Performance":  {"performance", PerfCommit},
		"Performances": {"performances", PerfCommit},

		"R":         {"r", RefactorCommit},
		"Refactor":  {"reFactor", RefactorCommit},
		"Refactors": {"reFactors", RefactorCommit},

		"S":      {"s", StyleCommit},
		"Style":  {"style", StyleCommit},
		"Styles": {"stYles", StyleCommit},

		"T":     {"t", TestCommit},
		"Test":  {"Test", TestCommit},
		"Tests": {"tests", TestCommit},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, FindCommitType(tc.str))
		})
	}
}

func TestParseCommitMsg(t *testing.T) {
	tcs := map[string]struct {
		str      string
		expected *CommitMessageOption
	}{
		"Bad": {"i'm bad", nil},
		"Simple": {"feat: message description",
			&CommitMessageOption{Ctype: FeatureCommit, Description: "message description"}},
		"Scoped": {"fix(scope): message description",
			&CommitMessageOption{Ctype: FixCommit, Description: "message description", Scope: "scope"}},
		"Breaking change": {"feat!: message description",
			&CommitMessageOption{Ctype: FeatureCommit, Description: "message description", BreakingChanges: true}},
		"With body": {"feat: message description\n\nCommit body\n",
			&CommitMessageOption{Ctype: FeatureCommit, Description: "message description", Body: "Commit body"}},
		"With footers": {"feat: message description\n\nCommit body\n\nFooter: 1\nFooter #2",
			&CommitMessageOption{Ctype: FeatureCommit, Description: "message description", Body: "Commit body", Footers: []string{"Footer: 1", "Footer #2"}}},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, ParseCommitMsg(tc.str))
		})
	}
}

func TestNextBump(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		cmsg string
		curr Bump
		next Bump
	}{
		{"Test next bump 1", "feat: a feature", BUMP_NONE, BUMP_MINOR},
		{"Test next bump 2", "feat: a feature", BUMP_PATCH, BUMP_MINOR},
		{"Test next bump 3", "feat: a feature", BUMP_MINOR, BUMP_MINOR},
		{"Test next bump 3", "feat: a feature", BUMP_MAJOR, BUMP_MAJOR},
		{"Test next bump 4", "fix: a fix", BUMP_NONE, BUMP_PATCH},
		{"Test next bump 5", "fix: a fix", BUMP_PATCH, BUMP_PATCH},
		{"Test next bump 6", "fix: a fix", BUMP_MINOR, BUMP_MINOR},
		{"Test next bump 7", "fix: a fix", BUMP_MAJOR, BUMP_MAJOR},
		{"Test next bump 8", "baadbeef", BUMP_NONE, BUMP_NONE},
		{"Test next bump 9", "baadbeef", BUMP_PATCH, BUMP_PATCH},
		{"Test next bump 10", "baadbeef", BUMP_MINOR, BUMP_MINOR},
		{"Test next bump 11", "baadbeef", BUMP_MAJOR, BUMP_MAJOR},
		{"Test next bump 12", "chore!: breaking", BUMP_NONE, BUMP_MAJOR},
		{"Test next bump 13", "chore!: breaking", BUMP_PATCH, BUMP_MAJOR},
		{"Test next bump 14", "chore!: breaking", BUMP_MINOR, BUMP_MAJOR},
		{"Test next bump 15", "chore!: breaking", BUMP_MAJOR, BUMP_MAJOR},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NextBump(tt.cmsg, tt.curr)
			assert.Equal(t, tt.next, actual)
		})
	}
}
