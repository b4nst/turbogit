package integrations

import (
	"os"
	"testing"

	"github.com/b4nst/turbogit/internal/format"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/nsf/termbox-go"
	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	id := IssueDescription{"ID-245", "Issue name", "An issue description.", "Jira", "type"}
	actual := id.Format(false)
	assert.Equal(t, "ID-245", id.ID)
	assert.Equal(t, "ID-245 - Issue name\n\nAn issue description.\n\nIssue provided by Jira", actual)
	actual = id.Format(true)
	assert.Contains(t, actual, "\x1B[0m")
}

func TestShortFormat(t *testing.T) {
	id := IssueDescription{"ID-245", "Issue name", "An issue description.", "Jira", "type"}
	actual := id.ShortFormat()
	assert.Equal(t, "ID-245 - Issue name", actual)
}

func TestSelectIssue(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping in CI due to termbox unstable behaviour.")
	}
	ids := []IssueDescription{
		{"ID-001", "Issue 1", "description", "Jira", "type"},
		{"ID-002", "Issue 2", "description", "Jira", "type"},
		{"ID-003", "Issue 3", "description", "Jira", "type"},
	}

	term := fuzzyfinder.UseMockedTerminal()
	term.SetSize(60, 10)
	term.SetEvents(
		termbox.Event{Type: termbox.EventKey, Ch: '2'},
		termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEnter},
		termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEsc})

	id, err := SelectIssue(ids, false)
	assert.NoError(t, err)
	assert.Equal(t, IssueDescription{"ID-002", "Issue 2", "description", "Jira", "type"}, id)
}

func TestToBranch(t *testing.T) {
	id := IssueDescription{Type: "feat", ID: "ID-245", Name: "feature 245."}
	expected := format.TugBranch{Type: id.Type, Prefix: id.ID, Description: id.Name}
	assert.Equal(t, expected, id.ToBranch(map[string]string{}))
}
