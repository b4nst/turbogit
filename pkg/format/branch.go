package format

import (
	"errors"
	"path"
	"regexp"
	"strings"
	"unicode"
)

var (
	forbiddenChar = regexp.MustCompile(`(?m)[\x60\?\*~^:\\\[\]]|@{|\.{2}`)
	blank         = regexp.MustCompile(`\s+`)
	void          = []byte("")
	sep           = []byte("-")

	// A default type rewrite map
	DefaultTypeRewrite = map[string]string{
		"feature": "feat",
		"bug":     "fix",
		"task":    "feat",
		"story":   "feat",
	}
)

// TugBranch represents a turbogit branch
type TugBranch struct {
	// Branch type (e.g. 'feat', 'fix', 'user', etc...)
	Type string
	// Branch prefix (issue id, user name, etc...)
	Prefix string
	// Branch description
	Description string
}

// String builds a git-sanitized branch name.
func (tb TugBranch) String() string {
	raw := path.Join(tb.Type, tb.Prefix, strings.ToLower(tb.Description))
	return sanitizeBranch(raw)
}

// ParseBranch parses a given string into a TugBranch or return an error on bad format.
func ParseBranch(s string) (TugBranch, error) {
	split := strings.SplitN(s, "/", 3)
	if len(split) < 2 {
		return TugBranch{}, errors.New("Bad branch format")
	}
	tb := TugBranch{}
	tb.Type = split[0]

	if len(split) < 3 {
		tb.Description = split[1]
	} else {
		tb.Prefix = split[1]
		tb.Description = split[2]
	}
	// Desanitize description
	desc := []rune(strings.ReplaceAll(tb.Description, "-", " "))
	desc[0] = unicode.ToUpper(desc[0])
	tb.Description = string(desc)

	return tb, nil
}

// WithType returns a TugBranch with the given type 't' or it's correlation in the rewrite map if it exists.
func (tb TugBranch) WithType(t string, rewrite map[string]string) TugBranch {
	ts := strings.ToLower(t)
	if tr, ok := rewrite[ts]; ok {
		tb.Type = tr
	} else {
		tb.Type = ts
	}

	return tb
}

func sanitizeBranch(s string) string {
	sb := forbiddenChar.ReplaceAll([]byte(s), void)
	s = string(blank.ReplaceAll(sb, sep))
	return strings.Trim(s, "./")
}
