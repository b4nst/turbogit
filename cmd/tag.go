/*
Copyright © 2020 banst

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"sort"
	"strings"

	"github.com/b4nst/turbogit/internal/format"
	"github.com/blang/semver/v4"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/spf13/cobra"
)

const (
	TAG_PREFIX = "v"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Create a tag",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("Not implemented")
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
}

type Tag struct {
	version semver.Version
	ref     *plumbing.Reference
}
type Tags []*Tag

func (slice Tags) Len() int {
	return len(slice)
}

func (slice Tags) Less(i, j int) bool {
	ti, tj := slice[i], slice[j]
	if ti == nil || tj == nil {
		return false
	}
	return ti.version.LT(tj.version)
}

func (slice Tags) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func filterSemver(it storer.ReferenceIter) (Tags, error) {
	tags := Tags{}

	filter := func(ref *plumbing.Reference) error {
		v, err := semver.Make(strings.TrimLeft(ref.Name().Short(), TAG_PREFIX))
		if err != nil {
			return err
		}
		tags = append(tags, &Tag{version: v, ref: ref})
		return nil
	}

	if err := it.ForEach(filter); err != nil {
		return nil, err
	}

	return tags, nil
}

// Return the last Semver tag or nil if
func lastTag(r *git.Repository) (*Tag, error) {
	iter, err := r.Tags()
	if err != nil {
		return nil, err
	}

	tags, err := filterSemver(iter)
	if err != nil {
		return nil, err
	}
	sort.Sort(sort.Reverse(tags))

	if len(tags) <= 0 {
		return nil, nil
	}
	return tags[0], nil
}

func nextVersion(curr semver.Version, msgs []string) (semver.Version, error) {
	const (
		Major int = iota
		Minor
		Patch
		Nil
	)
	next := Nil

	for _, msg := range msgs {
		cmo := format.ParseCommitMsg(msg)
		if cmo == nil {
			continue // Ignore malformatted commits
		}

		if cmo.BreakingChanges {
			next = Major
			break
		} else if cmo.Ctype == format.FeatureCommit {
			next = Minor
		} else if cmo.Ctype == format.FixCommit && next == Nil {
			next = Patch
		}
	}

	switch next {
	case Major:
		err := curr.IncrementMajor()
		return curr, err
	case Minor:
		err := curr.IncrementMinor()
		return curr, err
	case Patch:
		err := curr.IncrementPatch()
		return curr, err
	default:
		return curr, nil
	}
}
