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
	"fmt"
	"strings"
	"time"

	"github.com/b4nst/turbogit/internal/context"
	"github.com/b4nst/turbogit/internal/format"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().StringP("type", "t", "", "Commit type : [b]uild, [c]i, [d]ocs, f[e]at, [f]ix, [p]erf, [r]efactor, [s]tyle, [t]est")
	commitCmd.Flags().BoolP("breaking-changes", "c", false, "Commit contains breaking changes")
	commitCmd.Flags().StringP("body", "b", "", "Commit body.")
	commitCmd.Flags().StringP("scope", "s", "", "Commit scope.")
}

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit [subject]",
	Short: "Create a new commit.",
	Args:  cobra.MinimumNArgs(1),
	RunE:  commit,
}

func commit(cmd *cobra.Command, args []string) error {
	ctx, err := context.FromCommand(cmd)
	if err != nil {
		return err
	}

	msg := strings.Join(args, " ")

	ctype, err := validateCommitType(cmd)
	if err != nil {
		return err
	}
	if ctype == format.NilCommit {
		ctype, err = guessCommitType(msg, ctx)
		if err != nil {
			return err
		}
	}

	body, err := cmd.Flags().GetString("body")
	if err != nil {
		return err
	}

	bc, err := cmd.Flags().GetBool("breaking-changes")
	if err != nil {
		return err
	}

	scope, err := cmd.Flags().GetString("scope")
	if err != nil {
		return err
	}

	cmsg := format.CommitMessage(&format.CommitMessageOption{
		Ctype: ctype, Body: body, BreakingChanges: bc, Description: msg, Scope: scope,
	})

	// Commit
	w, err := ctx.Repo.Worktree()
	if err != nil {
		return err
	}

	author := object.Signature{
		Name:  ctx.Username,
		Email: ctx.Email,
		When:  time.Now(),
	}

	if _, err := w.Commit(cmsg, &git.CommitOptions{Author: &author}); err != nil {
		return err
	}

	return nil
}

func validateCommitType(cmd *cobra.Command) (format.CommitType, error) {
	ctype, err := cmd.Flags().GetString("type")
	if err != nil {
		return -1, err
	}

	// TODO refator that mess with a proper map
	switch ctype {
	case "":
		return format.NilCommit, nil
	case "b":
		return format.BuildCommit, nil
	case "c":
		return format.CiCommit, nil
	case "d":
		return format.DocCommit, nil
	case "e":
		return format.FeatureCommit, nil
	case "f":
		return format.FixCommit, nil
	case "p":
		return format.PerfCommit, nil
	case "r":
		return format.RefactorCommit, nil
	case "s":
		return format.StyleCommit, nil
	case "t":
		return format.TestCommit, nil
	default:
		return -1, fmt.Errorf("%s is not a commit type, allowed values are b, c, d, e, f, p, r, s, t", ctype)
	}
}

func guessCommitType(msg string, ctx *context.Context) (format.CommitType, error) {
	// TODO implement guessing feature
	return format.FeatureCommit, nil
}
