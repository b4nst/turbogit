package integrations

import (
	"fmt"
	"strings"

	"github.com/b4nst/turbogit/internal/format"
	"github.com/hpcloud/golor"
	"github.com/ktr0731/go-fuzzyfinder"
)

type IssueDescription struct {
	// Issue id
	ID string
	// Issue Name
	Name string
	// Issue description
	Description string
	// Issue provider
	Provider string
	// Issue type
	Type string
}

// Format format an issue description into a colored or raw string.
func (id IssueDescription) Format(color bool) string {
	var sb strings.Builder
	if color {
		id.ID = golor.Colorize(id.ID, golor.G, -1)
		id.Provider = golor.Colorize(id.Provider, golor.AssignColor(id.Provider), -1)
	}

	// Add ID
	sb.WriteString(id.ID)

	// ID - Title separator
	sb.WriteString(" - ")

	// Name
	sb.WriteString(id.Name)

	// Body
	if id.Description != "" {
		sb.WriteString("\n\n")
		sb.WriteString(id.Description)
	}

	// Provider
	sb.WriteString("\n\n")
	sb.WriteString("Issue provided by ")
	sb.WriteString(id.Provider)

	return sb.String()
}

// ShortFormat returns the shor representation string of an IssueDescription
func (id IssueDescription) ShortFormat() string {
	return fmt.Sprintf("%s - %s", id.ID, id.Name)
}

// ToBranch create a format.TugBranch from the issue description
func (id IssueDescription) ToBranch(rwtype map[string]string) format.TugBranch {
	return format.TugBranch{Prefix: id.ID, Description: id.Name}.WithType(id.Type, rwtype)
}

// SelectIssue prompts a fuzzy finder and returns the selected IssueDescription
// or an error if something unexpected happened
func SelectIssue(ids []IssueDescription, color bool) (IssueDescription, error) {
	idx, err := fuzzyfinder.Find(ids, func(i int) string {
		return ids[i].ShortFormat()
	},
		fuzzyfinder.WithPreviewWindow(func(i, _, _ int) string {
			if i == -1 {
				return ""
			}
			return ids[i].Format(color)
		}))
	if err != nil {
		return IssueDescription{}, err
	}
	return ids[idx], nil
}
