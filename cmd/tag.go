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
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/b4nst/turbogit/pkg/format"
	tugit "github.com/b4nst/turbogit/pkg/git"
	"github.com/blang/semver/v4"
	git "github.com/libgit2/git2go/v31"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(tagCmd)

	tagCmd.Flags().BoolP("dry-run", "d", false, "Do not tag.")
	tagCmd.Flags().StringP("prefix", "p", "v", "Tag prefix.")
}

type TagCmdOption struct {
	DryRun bool
	Prefix string
	Repo   *git.Repository
}

var tagCmd = &cobra.Command{
	Use:                   "tag",
	Short:                 "Create a tag",
	DisableFlagsInUseLine: true,
	Aliases:               []string{"release"},
	Long:                  "Create a semver tag, based on the commit history since last one",
	Example: `
# Given that the last release tag was v1.0.0, some feature were committed but no breaking changes.
# The following command will create the tag v1.1.0
$ tug tag
`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	Run:          runTagCmd,
}

func runTagCmd(cmd *cobra.Command, args []string) {
	tco, err := parseTagCmd(cmd, args)
	if err != nil {
		log.Fatal(err)
	}
	if err := runTag(tco); err != nil {
		log.Fatal(err)
	}
}

func parseTagCmd(cmd *cobra.Command, args []string) (*TagCmdOption, error) {
	// --dry-run
	fDryRun, err := cmd.Flags().GetBool("dry-run")
	if err != nil {
		return nil, err
	}

	// --prefix
	fPrefix, err := cmd.Flags().GetString("prefix")
	if err != nil {
		return nil, err
	}

	// Find repo
	repo, err := tugit.Getrepo()
	if err != nil {
		return nil, err
	}

	return &TagCmdOption{
		DryRun: fDryRun,
		Prefix: fPrefix,
		Repo:   repo,
	}, nil
}

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

func runTag(tco *TagCmdOption) error {
	r := tco.Repo

	walk, err := r.Walk()
	if err != nil {
		return err
	}
	if err := walk.PushHead(); err != nil {
		return err
	}

	bump := format.BUMP_NONE
	curr := semver.Version{}
	walker, err := commitWalker(&bump, &curr, tco.Prefix)
	if err != nil {
		return err
	}
	if err := walk.Iterate(walker); err != nil {
		return err
	}

	if bump == format.BUMP_NONE {
		fmt.Println("Nothing to do")
		return nil
	}
	// Bump tag
	if err := bumpVersion(&curr, bump); err != nil {
		return err
	}

	tagname := fmt.Sprintf("refs/tags/%s%s", tco.Prefix, curr)
	return tagHead(r, tagname, tco.DryRun)
}

func tagHead(r *git.Repository, tagname string, dry bool) error {
	head, err := r.Head()
	if err != nil {
		return err
	}
	if dry {
		fmt.Println(tagname, "would be created on", head.Target())
	} else {
		tag, err := r.References.Create(tagname, head.Target(), false, "")
		if err != nil {
			return err
		}
		fmt.Println(tag.Target(), "-->", tagname)
	}
	return nil
}

func bumpVersion(curr *semver.Version, bump format.Bump) error {
	if curr == nil {
		return errors.New("Received nil pointer for version")
	}
	switch bump {
	case format.BUMP_MAJOR:
		if curr.Major == 0 {
			return curr.IncrementMinor()
		}
		return curr.IncrementMajor()
	case format.BUMP_MINOR:
		return curr.IncrementMinor()
	case format.BUMP_PATCH:
		return curr.IncrementPatch()
	default:
		return nil
	}
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
