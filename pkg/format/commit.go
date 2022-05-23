package format

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/b4nst/turbogit/pkg/constants"
	"github.com/imdario/mergo"
)

type CommitType int

const (
	NilCommit CommitType = iota
	BuildCommit
	CiCommit
	ChoreCommit
	DocCommit
	FeatureCommit
	FixCommit
	PerfCommit
	RefactorCommit
	StyleCommit
	TestCommit
)

func (b CommitType) String() string {
	return [...]string{
		"",
		"build",
		"ci",
		"chore",
		"docs",
		"feat",
		"fix",
		"perf",
		"refactor",
		"style",
		"test",
	}[b]
}

func colorize(s string, code int) string {
	return fmt.Sprintf("\x1b[38;5;%03dm%s\x1b[0m", code, s)
}

func (b CommitType) ColorString() string {
	return [...]string{
		colorize("", 0),
		colorize("build", 200),
		colorize("ci", 92),
		colorize("chore", 15),
		colorize("docs", 250),
		colorize("feat", 2),
		colorize("fix", 1),
		colorize("perf", 3),
		colorize("refactor", 30),
		colorize("style", 6),
		colorize("test", 11),
	}[b]
}

func AllCommitType() []string {
	return []string{
		BuildCommit.String(),
		CiCommit.String(),
		ChoreCommit.String(),
		DocCommit.String(),
		FeatureCommit.String(),
		FixCommit.String(),
		PerfCommit.String(),
		RefactorCommit.String(),
		StyleCommit.String(),
		TestCommit.String(),
	}
}

var (
	buildCommitRe    = regexp.MustCompile(`(?i)^b(?:uilds?)?$`)
	ciCommitRe       = regexp.MustCompile(`(?i)^ci$`)
	choreCommitRe    = regexp.MustCompile(`(?i)^ch(?:ores?)?$`)
	docCommitRe      = regexp.MustCompile(`(?i)^d(?:ocs?)?$`)
	featureCommitRe  = regexp.MustCompile(`(?i)^fe(?:at(?:ure)?s?)?$`)
	fixCommitRe      = regexp.MustCompile(`(?i)^fi(?:x(?:es)?)?$`)
	perfCommitRe     = regexp.MustCompile(`(?i)^p(?:erf(:?ormance)?s?)?$`)
	refactorCommitRe = regexp.MustCompile(`(?i)^r(?:efactors?)?$`)
	styleCommitRe    = regexp.MustCompile(`(?i)^s(?:tyles?)?$`)
	testCommitRe     = regexp.MustCompile(`(?i)^t(?:ests?)?$`)
)

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

// Overwrite values with another CommitMessageOption.
func (cmo *CommitMessageOption) Overwrite(other *CommitMessageOption) error {
	return mergo.Merge(cmo, other, mergo.WithOverride)
}

// Check CommitMessageOption compliance
func (cmo *CommitMessageOption) Check() error {
	if cmo.Ctype == NilCommit {
		return errors.New("A commit type is required")
	}

	if cmo.Description == "" {
		return errors.New("A commit description is required")
	}

	return nil
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

// Extract type from string
func FindCommitType(str string) CommitType {
	s := []byte(str)
	switch {
	case buildCommitRe.Match(s):
		return BuildCommit
	case ciCommitRe.Match(s):
		return CiCommit
	case choreCommitRe.Match(s):
		return ChoreCommit
	case docCommitRe.Match(s):
		return DocCommit
	case featureCommitRe.Match(s):
		return FeatureCommit
	case fixCommitRe.Match(s):
		return FixCommit
	case perfCommitRe.Match(s):
		return PerfCommit
	case refactorCommitRe.Match(s):
		return RefactorCommit
	case styleCommitRe.Match(s):
		return StyleCommit
	case testCommitRe.Match(s):
		return TestCommit
	default:
		return NilCommit
	}
}

func ParseCommitMsg(msg string) *CommitMessageOption {
	lines := strings.Split(msg, "\n")

	// First line
	re := regexp.MustCompile(`(?m)^(?P<type>\w+)(?:\((?P<scope>[^)]+)\))?(?P<bc>!)?: (?P<subject>.+)$`)
	match := re.FindStringSubmatch(lines[0])
	if len(match) <= 0 {
		return nil
	}
	res := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			res[name] = match[i]
		}
	}
	cmo := &CommitMessageOption{
		Ctype:           FindCommitType(res["type"]),
		Description:     res["subject"],
		Scope:           res["scope"],
		BreakingChanges: res["bc"] == "!",
	}

	// Body and footers
	re = regexp.MustCompile(`(?m)^\w+(?: #|: )`)
	for _, l := range lines[1:] {
		if re.MatchString(l) {
			cmo.Footers = append(cmo.Footers, l)
		} else {
			cmo.Body += l
		}
	}
	cmo.Body = strings.Trim(cmo.Body, "\n")

	return cmo
}

type Bump int

const (
	BUMP_NONE Bump = iota
	BUMP_PATCH
	BUMP_MINOR
	BUMP_MAJOR
)

// Get the next bump that this commit message would generate
func NextBump(cmsg string, curr Bump) Bump {
	if curr == BUMP_MAJOR {
		return curr
	}
	co := ParseCommitMsg(cmsg)
	if co == nil {
		return curr
	}
	if co.BreakingChanges {
		return BUMP_MAJOR
	}
	if co.Ctype == FeatureCommit {
		return BUMP_MINOR
	}
	if curr < BUMP_PATCH && co.Ctype == FixCommit {
		return BUMP_PATCH
	}

	return curr
}
