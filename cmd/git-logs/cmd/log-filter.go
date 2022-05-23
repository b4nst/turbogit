package cmd

import (
	"time"

	"github.com/b4nst/turbogit/pkg/format"
	git "github.com/libgit2/git2go/v33"
)

type LogFilter func(c *git.Commit, co *format.CommitMessageOption) (keep, walk bool)

var PassThru LogFilter = func(c *git.Commit, co *format.CommitMessageOption) (bool, bool) { return true, true }

func ApplyFilters(c *git.Commit, co *format.CommitMessageOption, filters ...LogFilter) (keep, walk bool) {
	for _, filter := range filters {
		keep, walk = filter(c, co)
		if !walk || !keep {
			return
		}
	}
	return true, true
}

func Since(since *time.Time) LogFilter {
	if since == nil {
		return PassThru
	}

	return func(c *git.Commit, co *format.CommitMessageOption) (keep, walk bool) {
		d := c.Committer().When
		if d.Before(*since) {
			return false, false
		}
		return true, true
	}
}

func Until(until *time.Time) LogFilter {
	if until == nil {
		return PassThru
	}

	return func(c *git.Commit, co *format.CommitMessageOption) (keep, walk bool) {
		d := c.Committer().When
		if d.After(*until) {
			return false, true
		}
		return true, true
	}
}

func BreakingChange(is bool) LogFilter {
	return func(c *git.Commit, co *format.CommitMessageOption) (keep, walk bool) {
		return co.BreakingChanges == is, true
	}
}

func Type(types []format.CommitType) LogFilter {
	if len(types) <= 0 {
		return PassThru
	}

	mt := make(map[format.CommitType]bool, len(types))
	for _, ct := range types {
		mt[ct] = true
	}
	return func(c *git.Commit, co *format.CommitMessageOption) (keep, walk bool) {
		return mt[co.Ctype], true
	}
}

func Scope(scopes []string) LogFilter {
	if len(scopes) <= 0 {
		return PassThru
	}

	ms := make(map[string]bool, len(scopes))
	for _, s := range scopes {
		ms[s] = true
	}
	return func(c *git.Commit, co *format.CommitMessageOption) (keep, walk bool) {
		return ms[co.Scope], true
	}
}
