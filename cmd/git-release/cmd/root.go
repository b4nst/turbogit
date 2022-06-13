/*
Copyright © 2022 banst

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

	"github.com/b4nst/turbogit/internal/cmdbuilder"
	"github.com/b4nst/turbogit/pkg/format"
	"github.com/blang/semver/v4"
	git "github.com/libgit2/git2go/v33"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.Flags().BoolP("dry-run", "d", false, "Do not tag.")
	RootCmd.Flags().StringP("prefix", "p", "v", "Tag prefix.")

	cmdbuilder.RepoAware(RootCmd)
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "git-release",
	Short: "Release a SemVer tag based on the commit history.",
	Example: `
# Given that the last release tag was v1.0.0, some feature were committed but no breaking changes.
# The following command will create the tag v1.1.0
$ git release
`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	Run:          run,
}

func run(cmd *cobra.Command, args []string) {
	// get options
	dryrun, err := cmd.Flags().GetBool("dry-run")
	cobra.CheckErr(err)
	prefix, err := cmd.Flags().GetString("prefix")
	cobra.CheckErr(err)
	r := cmdbuilder.GetRepo(cmd)

	// initialize walker
	walk, err := r.Walk()
	cobra.CheckErr(err)
	cobra.CheckErr(walk.PushHead())

	// find next version
	bump := format.BUMP_NONE
	curr := semver.Version{}
	walker, err := commitWalker(&bump, &curr, prefix)
	cobra.CheckErr(err)
	cobra.CheckErr(walk.Iterate(walker))

	if bump == format.BUMP_NONE {
		fmt.Println("Nothing to do")
		return
	}
	// Bump tag
	cobra.CheckErr(bumpVersion(&curr, bump))

	// do tag
	tagname := fmt.Sprintf("refs/tags/%s%s", prefix, curr)
	cobra.CheckErr(tagHead(r, tagname, dryrun))
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
		return errors.New("current version must not be nil")
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
