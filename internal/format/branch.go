package format

import (
	"fmt"
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
	return [...]string{"feat", "fix", "user"}[b]
}

func AllBranchType() []string {
	return []string{
		FeatureBranch.String(),
		FixBranch.String(),
		UserBranch.String(),
	}
}

func BranchTypeFrom(str string) (BranchType, error) {
	switch str {
	case FeatureBranch.String():
		return FeatureBranch, nil
	case FixBranch.String():
		return FixBranch, nil
	case UserBranch.String():
		return UserBranch, nil
	default:
		return -1, fmt.Errorf("%s is not a branch type, allowed values are f,p and u", str)
	}
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
