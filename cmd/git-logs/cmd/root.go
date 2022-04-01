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
	"fmt"
	"io"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/araddon/dateparse"
	"github.com/b4nst/turbogit/pkg/format"
	tugit "github.com/b4nst/turbogit/pkg/git"
	"github.com/hpcloud/golor"
	git "github.com/libgit2/git2go/v33"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.Flags().BoolP("all", "a", false, "Pretend as if all the refs in refs/, along with HEAD, are listed on the command line as <commit>. If set on true, the --from option will be ignored.")
	rootCmd.Flags().Bool("no-color", false, "Disable color output")
	rootCmd.Flags().StringP("from", "f", "HEAD", "Logs only commits reachable from this one")
	rootCmd.Flags().String("since", "", "Show commits more recent than a specific date")
	rootCmd.Flags().String("until", "", "Show commits older than a specific date")
	// logCmd.Flags().String("path", "", "Filter commits based on the path of files that are updated. Accept regexp")
	// Filters
	rootCmd.Flags().StringArrayP("type", "t", []string{}, "Filter commits by type (repeatable option)")
	rootCmd.Flags().StringArrayP("scope", "s", []string{}, "Filter commits by scope (repeatable option)")
	rootCmd.Flags().BoolP("breaking-changes", "c", false, "Only shows breaking changes")
}

type option struct {
	All            bool
	NoColor        bool
	From           string
	Since          *time.Time
	Until          *time.Time
	Types          []format.CommitType
	Scopes         []string
	BreakingChange bool
	Repo           *git.Repository
}

// rootCmd represents the log command
var rootCmd = &cobra.Command{
	Use:   "git-logs",
	Short: "Shows the commit logs.",
	Args:  cobra.NoArgs,
	Run:   runCmd,
}

func runCmd(cmd *cobra.Command, args []string) {
	opt, err := parseCmd(cmd, args)
	if err != nil {
		log.Fatal(err)
	}
	if err := run(opt); err != nil {
		log.Fatal(err)
	}
}

func parseCmd(cmd *cobra.Command, args []string) (*option, error) {
	// --all
	fAll, err := cmd.Flags().GetBool("all")
	if err != nil {
		return nil, err
	}
	// --no-color
	fNoColor, err := cmd.Flags().GetBool("no-color")
	if err != nil {
		return nil, err
	}
	// --from
	fFrom, err := cmd.Flags().GetString("from")
	if err != nil {
		return nil, err
	}
	// --since
	fSince, err := cmd.Flags().GetString("since")
	if err != nil {
		return nil, err
	}
	var since *time.Time
	if fSince != "" {
		date, err := dateparse.ParseAny(fSince)
		if err != nil {
			return nil, err
		}
		since = &date
	}
	// --until
	fUntil, err := cmd.Flags().GetString("until")
	if err != nil {
		return nil, err
	}
	var until *time.Time
	if fUntil != "" {
		date, err := dateparse.ParseAny(fUntil)
		if err != nil {
			return nil, err
		}
		until = &date
	}
	// --types
	fTypes, err := cmd.Flags().GetStringArray("type")
	if err != nil {
		return nil, err

	}
	types := make([]format.CommitType, len(fTypes))
	for i, v := range fTypes {
		types[i] = format.FindCommitType(v)
		// TODO warn or error on nil commit type
	}
	// --scopes
	fScopes, err := cmd.Flags().GetStringArray("scope")
	if err != nil {
		return nil, err
	}
	// --breaking-changes
	fBreakingChanges, err := cmd.Flags().GetBool("breaking-changes")
	if err != nil {
		return nil, err
	}

	// Find repo
	repo, err := tugit.Getrepo()
	if err != nil {
		return nil, err
	}

	return &option{
		All:            fAll,
		NoColor:        fNoColor,
		From:           fFrom,
		Since:          since,
		Until:          until,
		Types:          types,
		Scopes:         fScopes,
		BreakingChange: fBreakingChanges,
		Repo:           repo,
	}, nil
}

func run(opt *option) error {
	r := opt.Repo

	walk, err := r.Walk()
	if err != nil {
		return err
	}
	if opt.All {
		if err := walk.PushGlob("refs/*"); err != nil {
			return err
		}
	} else {
		from, err := r.RevparseSingle(opt.From)
		if err != nil {
			return err
		}
		if err := walk.Push(from.Id()); err != nil {
			return err
		}
	}

	// Build filters
	filters := []LogFilter{}
	if opt.Since != nil {
		filters = append(filters, func(c *git.Commit, co *format.CommitMessageOption) (p, continueWalk bool) {
			d := c.Committer().When
			if d.Before(*opt.Since) {
				return false, false
			}
			return true, true
		})
	}
	if opt.Until != nil {
		filters = append(filters, func(c *git.Commit, co *format.CommitMessageOption) (p, continueWalk bool) {
			d := c.Committer().When
			if d.After(*opt.Until) {
				return false, true
			}
			return true, true
		})
	}
	if opt.BreakingChange {
		filters = append(filters, func(c *git.Commit, co *format.CommitMessageOption) (p, continueWalk bool) {
			return co.BreakingChanges, true
		})
	}
	if len(opt.Types) > 0 {
		filters = append(filters, func(c *git.Commit, co *format.CommitMessageOption) (p, continueWalk bool) {
			for _, t := range opt.Types {
				if co.Ctype == t {
					return true, true
				}
			}
			return false, true
		})
	}
	if len(opt.Scopes) > 0 {
		filters = append(filters, func(c *git.Commit, co *format.CommitMessageOption) (p, continueWalk bool) {
			for _, s := range opt.Scopes {
				if co.Scope == s {
					return true, true
				}
			}
			return false, true
		})
	}
	// Writer
	w := tabwriter.NewWriter(os.Stdout, 8, 8, 0, ' ', 0)

	if err := walk.Iterate(buildLogWalker(w, !opt.NoColor, filters)); err != nil {
		return err
	}

	return nil
}

type LogFilter func(c *git.Commit, co *format.CommitMessageOption) (p, continueWalk bool)

func buildLogWalker(w io.Writer, color bool, filters []LogFilter) func(c *git.Commit) bool {
	return func(c *git.Commit) bool {
		co := format.ParseCommitMsg(c.Message())
		parsed := true
		if co == nil {
			parsed = false
			co = &format.CommitMessageOption{}
		}
		p, continueWalk := true, true
		for _, filter := range filters {
			p, continueWalk = filter(c, co)
			if !continueWalk {
				return false
			}
			if !p {
				break
			}
		}
		if p {
			fprettyprint(w, c, co, color, parsed)
		}
		return true
	}
}

func fprettyprint(w io.Writer, c *git.Commit, co *format.CommitMessageOption, color bool, parsed bool) {
	// Hash
	h, err := c.ShortId()
	if err != nil {
		h = c.Id().String()
	}
	if color {
		h = golor.Colorize(h, golor.W, -1)
	}
	fmt.Fprintf(w, "(%s)", h)

	// Message
	msg := co.Description
	if color {
		msg = golor.Colorize(msg, 15, -1)
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
	author := c.Author()
	fmt.Fprintf(w, "\tAuthor:\t%s <%s>\n", author.Name, author.Email)
	// Committer
	committer := c.Committer()
	fmt.Fprintf(w, "\tCommitter:\t%s <%s>\n", committer.Name, committer.Email)

	// Date
	fmt.Fprintf(w, "\tDate:\t%s\n", committer.When.Format(time.UnixDate))

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
