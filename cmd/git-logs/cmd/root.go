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
	filters := []LogFilter{
		Since(opt.Since),
		Until(opt.Until),
		Type(opt.Types),
		Scope(opt.Scopes),
		BreakingChange(opt.BreakingChange),
	}

	tw := tabwriter.NewWriter(os.Stdout, 10, 1, 1, ' ', 0)
	defer tw.Flush()
	if err := walk.Iterate(buildLogWalker(tw, !opt.NoColor, filters)); err != nil {
		return err
	}

	return nil
}

func buildLogWalker(w io.Writer, color bool, filters []LogFilter) func(c *git.Commit) bool {
	return func(c *git.Commit) bool {
		co := format.ParseCommitMsg(c.Message())
		if co == nil {
			co = &format.CommitMessageOption{}
		}
		keep, walk := ApplyFilters(c, co, filters...)
		if !keep {
			return walk
		}

		// Hash
		h, err := c.ShortId()
		if err != nil {
			h = c.Id().String()
		}
		// type
		var ctype string
		if color {
			h = fmt.Sprintf("\x1b[38;5;231m%s\x1b[0m", h)
			ctype = co.Ctype.ColorString()
		} else {
			ctype = co.Ctype.String()
		}
		// description
		msg := co.Description
		if msg == "" {
			msg = c.Summary()
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t\n", h, ctype, msg)

		return walk
	}
}
