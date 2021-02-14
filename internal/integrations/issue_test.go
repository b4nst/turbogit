package integrations

import (
	"testing"

	"github.com/b4nst/turbogit/internal/format"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/nsf/termbox-go"
	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	id := IssueDescription{"ID-245", "Issue name", "An issue description.", "Jira", "type"}
	actual := id.Format(true)
	assert.Equal(t, "ID-245", id.ID)
	assert.Equal(t, "\x1b[1;32;89mID-245\x1b[0m - Issue name\n\nAn issue description.\n\nIssue provided by \x1b[0;36;89mJira\x1b[0m", actual)
}

func TestShortFormat(t *testing.T) {
	id := IssueDescription{"ID-245", "Issue name", "An issue description.", "Jira", "type"}
	actual := id.ShortFormat()
	assert.Equal(t, "ID-245 - Issue name", actual)
}

func TestSelectIssue(t *testing.T) {
	ids := []IssueDescription{
		{"ID-001", "Issue 1", "description", "Jira", "type"},
		{"ID-002", "Issue 2", "description", "Jira", "type"},
		{"ID-003", "Issue 3", "description", "Jira", "type"},
	}

	keys := func(str string) []termbox.Event {
		s := []rune(str)
		e := make([]termbox.Event, 0, len(s))
		for _, r := range s {
			e = append(e, termbox.Event{Type: termbox.EventKey, Ch: r})
		}
		return e
	}

	term := fuzzyfinder.UseMockedTerminal()
	term.SetSize(60, 10)
	term.SetEvents(append(
		keys("issue 2"),
		termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEnter})...)

	id, err := SelectIssue(ids, false)
	assert.NoError(t, err)
	assert.Equal(t, IssueDescription{"ID-002", "Issue 2", "description", "Jira", "type"}, id)
}

func TestToBranch(t *testing.T) {
	id := IssueDescription{Type: "feat", ID: "ID-245", Name: "feature 245."}
	expected := format.TugBranch{Type: id.Type, Prefix: id.ID, Description: id.Name}
	assert.Equal(t, expected, id.ToBranch(map[string]string{}))
}
