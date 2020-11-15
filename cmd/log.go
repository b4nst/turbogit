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
	"strings"
	"text/tabwriter"

	"github.com/b4nst/turbogit/internal/context"
	"github.com/b4nst/turbogit/internal/format"
	"github.com/go-git/go-git/v5"
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
}

func runLog(cmd *cobra.Command, args []string) error {
	ctx, err := context.FromCommand(cmd)
	if err != nil {
		return err
	}

	cIter, err := ctx.Repo.Log(&git.LogOptions{Order: git.LogOrderCommitterTime})
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 8, 8, 0, ' ', 0)
	cIter.ForEach(func(c *object.Commit) error {
		fprettyprint(w, c, 7, true)
		fmt.Fprintln(w)
		return nil
	})
	return nil
}

// Format commit to pretty string (mostly for log usage)
func fprettyprint(w io.Writer, c *object.Commit, hLength int, color bool) {
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

	// Hash
	h := c.Hash.String()
	if hLength > len(h) || hLength < 0 {
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
