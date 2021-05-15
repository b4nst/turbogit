package integrations

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andygrunwald/go-jira"
	git "github.com/libgit2/git2go/v30"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJiraProvider(t *testing.T) {
	cf, err := ioutil.TempFile("", "gitconfig")
	require.NoError(t, err)
	c, err := git.OpenOndisk(cf.Name())
	require.NoError(t, err)

	p, err := jiraProvider(c)
	assert.NoError(t, err)
	assert.Nil(t, p)

	// No username
	err = c.SetBool("jira.enable", true)
	require.NoError(t, err)
	p, err = jiraProvider(c)
	assert.EqualError(t, err, "config value 'jira.username' was not found")
	assert.Nil(t, p)

	// No token
	err = c.SetString("jira.username", "alice@ecorp.com")
	require.NoError(t, err)
	p, err = jiraProvider(c)
	assert.EqualError(t, err, "config value 'jira.token' was not found")
	assert.Nil(t, p)

	// No domain
	err = c.SetString("jira.token", "supersecret")
	require.NoError(t, err)
	p, err = jiraProvider(c)
	assert.EqualError(t, err, "config value 'jira.domain' was not found")
	assert.Nil(t, p)

	// No filter
	err = c.SetString("jira.domain", "foo.bar")
	require.NoError(t, err)
	p, err = jiraProvider(c)
	assert.EqualError(t, err, "config value 'jira.filter' was not found")
	assert.Nil(t, p)

	// All's ok
	err = c.SetString("jira.filter", "query filter")
	require.NoError(t, err)
	p, err = jiraProvider(c)
	assert.NoError(t, err)
	assert.IsType(t, &JiraProvider{}, p)
}

func TestSearch(t *testing.T) {
	filter := "foofilter"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, filter, r.URL.Query().Get("jql"))

		js, err := json.Marshal(struct {
			Issues     []jira.Issue `json:"issues" structs:"issues"`
			StartAt    int          `json:"startAt" structs:"startAt"`
			MaxResults int          `json:"maxResults" structs:"maxResults"`
			Total      int          `json:"total" structs:"total"`
		}{
			StartAt:    0,
			MaxResults: 50,
			Total:      1,
			Issues: []jira.Issue{
				{
					Key: "B#NST",
					Fields: &jira.IssueFields{
						Summary:     "issue",
						Description: "description",
						Type: jira.IssueType{
							Name: "type",
						},
					},
				},
			},
		})
		require.NoError(t, err)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}))
	defer ts.Close()

	client, err := jira.NewClient(ts.Client(), ts.URL)
	require.NoError(t, err)
	provider := JiraProvider{
		filter: filter,
		client: client,
	}
	ids, err := provider.Search()
	assert.NoError(t, err)
	assert.Len(t, ids, 1)
	assert.Equal(t, IssueDescription{
		ID:          "B#NST",
		Name:        "issue",
		Description: "description",
		Type:        "type",
		Provider:    JIRA_PROVIDER,
	}, ids[0])

}
