package cmd

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/b4nst/turbogit/pkg/format"
	"github.com/blang/semver/v4"
	git "github.com/libgit2/git2go/v33"
)

func commitWalker(bump *format.Bump, curr *semver.Version, prefix string) (func(*git.Commit) bool, error) {
	dfo, err := git.DefaultDescribeFormatOptions()
	if err != nil {
		return nil, err
	}
	dco := &git.DescribeOptions{
		MaxCandidatesTags:     1,
		Strategy:              git.DescribeTags,
		Pattern:               fmt.Sprintf("%s*", prefix),
		OnlyFollowFirstParent: true,
	}

	return func(c *git.Commit) bool {
		dr, err := c.Describe(dco)
		if err != nil {
			// No next tag matching
			*bump = format.NextBump(c.Message(), *bump)
			return true
		}
		d, err := dr.Format(&dfo)
		if err != nil {
			panic(err)
		}
		v, offset, err := parseDescription(d)
		*curr = *v
		if err != nil {
			panic(err)
		}
		if offset <= 1 {
			return false
		}
		*bump = format.NextBump(c.Message(), *bump)
		return true
	}, nil
}

func parseDescription(d string) (*semver.Version, int, error) {
	re, err := regexp.Compile(`-(\d+)-[a-z0-9]{8}$`)
	if err != nil {
		return nil, 0, err
	}
	offset := 1

	if res := re.FindStringSubmatch(d); res != nil {
		offset, err = strconv.Atoi(res[1])
		if err != nil {
			return nil, 0, err
		}
		offset++
		d = d[:len(d)-len(res[0])]
	}

	if v, err := semver.ParseTolerant(d); err == nil {
		return &v, offset, nil
	}
	return nil, offset, nil
}
