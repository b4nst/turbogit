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
	void          = []byte("")
)

func sanitize(s string) string {
	s = string(forbiddenChar.ReplaceAll([]byte(s), void))
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.Trim(s, "/")
	return strings.ToLower(s)
}

func BranchName(btype BranchType, description string, username string) string {
	branch := btype.String()
	if btype == UserBranch {
		branch += "/" + username
	}
	if description != "" {
		branch += "/" + sanitize(description)
	}

	return branch
}
