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
	"io"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/araddon/dateparse"
	"github.com/b4nst/turbogit/internal/context"
	"github.com/b4nst/turbogit/internal/format"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hpcloud/golor"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:                   "log",
	Short:                 "Shows the commit logs.",
	DisableFlagsInUseLine: true,
	SilenceUsage:          true,
	Args:                  cobra.NoArgs,
	RunE:                  runLog,
}

func init() {
	RootCmd.AddCommand(logCmd)

	logCmd.Flags().BoolP("all", "a", false, "Pretend as if all the refs in refs/, along with HEAD, are listed on the command line as <commit>. If set on true, the --from option will be ignored.")
	logCmd.Flags().Bool("no-color", false, "Disable color output")
	logCmd.Flags().StringP("from", "f", "HEAD", "Logs only commits reachable from this one")
	logCmd.Flags().String("since", "", "Show commits more recent than a specific date")
	logCmd.Flags().String("until", "", "Show commits older than a specific date")
	logCmd.Flags().String("path", "", "Filter commits based on the path of files that are updated. Accept regexp")
	logCmd.Flags().Int("hash-length", 7, "Commit hash length. Set a value <=0 to get the full length")
	// Filters
	logCmd.Flags().StringArrayP("type", "t", []string{}, "Filter commits by type (repeatable option)")
	commitCmd.RegisterFlagCompletionFunc("type", typeFlagCompletion)
	logCmd.Flags().StringArrayP("scope", "s", []string{}, "Filter commits by scope (repeatable option)")
	logCmd.Flags().BoolP("breaking-changes", "c", false, "Only shows breaking changes")
}

func runLog(cmd *cobra.Command, args []string) error {
	ctx, err := context.FromCommand(cmd)
	if err != nil {
		return err
	}

	// --all
	fAll, err := cmd.Flags().GetBool("all")
	if err != nil {
		return err
	}
	// --no-color
	fNoColor, err := cmd.Flags().GetBool("no-color")
	if err != nil {
		return err
	}
	// --from
	fFrom, err := cmd.Flags().GetString("from")
	if err != nil {
		return err
	}
	// --hash-length
	fHLength, err := cmd.Flags().GetInt("hash-length")
	if err != nil {
		return err
	}
	// --since
	fSince, err := cmd.Flags().GetString("since")
	if err != nil {
		return err
	}
	var since *time.Time
	if fSince != "" {
		*since, err = dateparse.ParseAny(fSince)
		if err != nil {
			return err
		}
	}
	// --until
	fUntil, err := cmd.Flags().GetString("until")
	if err != nil {
		return err
	}
	var until *time.Time
	if fUntil != "" {
		*until, err = dateparse.ParseAny(fUntil)
		if err != nil {
			return err
		}
	}
	// --path
	fPath, err := cmd.Flags().GetString("path")
	if err != nil {
		return err
	}
	var pathFilter func(string) bool
	if fPath != "" {
		pathReg, err := regexp.Compile(fPath)
		if err != nil {
			return err
		}
		pathFilter = func(p string) bool { return pathReg.MatchString(p) }
	}
	// Filters
	filters := []FilterCommit{}
	fTypes, err := cmd.Flags().GetStringArray("type")
	if err != nil {
		return err
	}
	if len(fTypes) > 0 {
		types := make([]format.CommitType, len(fTypes))
		for i, v := range fTypes {
			types[i] = format.FindCommitType(v)
		}
		filters = append(filters, func(c *object.Commit, co *format.CommitMessageOption) bool {
			for _, t := range types {
				if co.Ctype == t {
					return true
				}
			}
			return false
		})
	}
	fScopes, err := cmd.Flags().GetStringArray("scope")
	if err != nil {
		return err
	}
	if len(fScopes) > 0 {
		filters = append(filters, func(c *object.Commit, co *format.CommitMessageOption) bool {
			for _, s := range fScopes {
				if co.Scope == s {
					return true
				}
			}
			return false
		})
	}
	fBreakingChanges, err := cmd.Flags().GetBool("breaking-changes")
	if err != nil {
		return err
	}
	if fBreakingChanges {
		filters = append(filters, func(c *object.Commit, co *format.CommitMessageOption) bool {
			return co.BreakingChanges
		})
	}

	from, err := ctx.Repo.ResolveRevision(plumbing.Revision(fFrom))
	if err != nil {
		return fmt.Errorf("Error looking for commit %s: %s", fFrom, err)
	}
	cIter, err := ctx.Repo.Log(&git.LogOptions{Order: git.LogOrderCommitterTime, All: fAll, From: *from, Since: since, Until: until, PathFilter: pathFilter})
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 8, 8, 0, ' ', 0)
	cIter.ForEach(func(c *object.Commit) error {
		co := format.ParseCommitMsg(c.Message)
		parsed := true
		if co == nil {
			parsed = false
			sm := strings.SplitN(c.Message, "\n", 2)
			co = &format.CommitMessageOption{
				Description: sm[0],
				Ctype:       format.NilCommit,
			}
			if len(sm) > 1 {
				co.Body = strings.Join(sm[1:], "\n")
			}
		}
		if applyFilters(c, co, filters...) {
			fprettyprint(w, c, co, fHLength, !fNoColor, parsed)
			fmt.Fprintln(w)
		}
		return nil
	})
	return nil
}

type FilterCommit func(*object.Commit, *format.CommitMessageOption) bool

func applyFilters(c *object.Commit, co *format.CommitMessageOption, filters ...FilterCommit) bool {
	keep := true
	for _, f := range filters {
		if !f(c, co) {
			keep = false
			break
		}
	}
	return keep
}

func fprettyprint(w io.Writer, c *object.Commit, co *format.CommitMessageOption, hLength int, color bool, parsed bool) {
	// Hash
	h := c.Hash.String()
	if hLength > len(h) || hLength <= 0 {
		hLength = len(h)
	}
	h = h[0:hLength] // Truncate
	if color {
		h = golor.Colorize(h, golor.W, -1)
	}
	fmt.Fprintf(w, "(%s)", h)

	// Message
	msg := co.Description
	if color {
		msg = golor.Colorize(msg, 215, -1)
	}
	fmt.Fprintf(w, " %s", msg)

	// Annotation
	if co.BreakingChanges || !parsed {
		an := "!BREAKING CHANGE"
		if !parsed {
			an = "!BADBEEF"
		}
		if color {
			an = golor.Colorize(an, golor.W, golor.RED)
		}
		fmt.Fprintf(w, " - %s", an)
	}
	// End of the first line
	fmt.Fprintln(w)

	// Author
	author := c.Author.String()
	if color {
		author = golor.Colorize(author, golor.AssignColor(author), -1)
	}
	fmt.Fprintf(w, "\tAuthor:\t%s\n", author)

	// Date
	fmt.Fprintf(w, "\tDate:\t%s\n", c.Author.When.Format(object.DateFormat))

	if parsed {
		// Type
		ct := co.Ctype.String()
		if color {
			ct = format.ColorizeCommitType(ct, co.Ctype)
		}
		fmt.Fprintf(w, "\tType:\t%s\n", ct)
		// Scope
		scope := co.Scope
		if scope == "" {
			scope = "none"
		}
		if color {
			scope = golor.Colorize(scope, golor.AssignColor(scope), -1)
		}
		fmt.Fprintf(w, "\tScope:\t%s\n", scope)
	}

	if co.Body != "" || len(co.Footers) > 0 {
		fmt.Fprintf(w, "\n\t%s", co.Body)
		for _, f := range co.Footers {
			fmt.Fprintf(w, "\n\t%s", f)
		}
		fmt.Fprintln(w)
	}
}
