package format

import (
	"regexp"
	"strings"
)

type BranchType int

const (
	FeatureBranch BranchType = iota
	FixBranch
	UserBranch
)

func (b BranchType) String() string {
	return [...]string{"features", "fix", "users"}[b]
}

var (
	forbiddenChar = regexp.MustCompile(`(?m)[\?\*~^:\\]|@{|\.{2}`)
	blank         = regexp.MustCompile(`\s+`)
	void          = []byte("")
	sep           = []byte("-")
)

func sanitizeBranch(s string) string {
	sb := forbiddenChar.ReplaceAll([]byte(s), void)
	s = string(blank.ReplaceAll(sb, sep))
	s = strings.Trim(s, "/")
	return strings.ToLower(s)
}

func BranchName(btype BranchType, description string, username string) string {
	branch := btype.String()
	if btype == UserBranch {
		branch += "/" + username
	}
	if description != "" {
		branch += "/" + description
	}

	return sanitizeBranch(branch)
}
