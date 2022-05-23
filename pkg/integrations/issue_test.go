package integrations

import (
	"testing"

	"github.com/b4nst/turbogit/pkg/format"
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

func TestToBranch(t *testing.T) {
	id := IssueDescription{Type: "feat", ID: "ID-245", Name: "feature 245."}
	expected := format.TugBranch{Type: id.Type, Prefix: id.ID, Description: id.Name}
	assert.Equal(t, expected, id.ToBranch(map[string]string{}))
}
