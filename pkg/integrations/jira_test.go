package integrations

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andygrunwald/go-jira"
	"github.com/b4nst/turbogit/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJiraProvider(t *testing.T) {
	r := test.TestRepo(t)
	defer test.CleanupRepo(t, r)
	test.InitRepoConf(t, r)

	p, err := NewJiraProvider(r)
	assert.NoError(t, err)
	assert.Nil(t, p)

	c, err := r.Config()
	require.NoError(t, err)
	require.NoError(t, c.SetBool("jira.enabled", true))

	// No username
	p, err = NewJiraProvider(r)
	assert.EqualError(t, err, "config value 'jira.username' was not found")
	assert.Nil(t, p)

	// No token
	require.NoError(t, c.SetString("jira.username", "alice@ecorp.com"))
	p, err = NewJiraProvider(r)
	assert.EqualError(t, err, "config value 'jira.token' was not found")
	assert.Nil(t, p)

	// No domain
	require.NoError(t, c.SetString("jira.token", "supersecret"))
	p, err = NewJiraProvider(r)
	assert.EqualError(t, err, "config value 'jira.domain' was not found")
	assert.Nil(t, p)

	// No filter
	require.NoError(t, c.SetString("jira.domain", "foo.bar"))
	p, err = NewJiraProvider(r)
	assert.EqualError(t, err, "config value 'jira.filter' was not found")
	assert.Nil(t, p)

	// All's ok
	require.NoError(t, c.SetString("jira.filter", "query filter"))
	p, err = NewJiraProvider(r)
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
