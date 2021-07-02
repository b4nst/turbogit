package integrations

import (
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/briandowns/spinner"
	git "github.com/libgit2/git2go/v31"
)

const (
	// Jira provider's name
	JIRA_PROVIDER = "Jira"
)

// JiraProvider represents the Jira issue provider.
type JiraProvider struct {
	filter string
	client *jira.Client
}

// Search returns a list of issues matching the query or an error if the request failed.
func (jp JiraProvider) Search() ([]IssueDescription, error) {
	sopts := &jira.SearchOptions{}

	s := spinner.New(spinner.CharSets[39], 100*time.Millisecond)
	s.Suffix = " Searching on Jira"
	s.Start()
	raw, _, err := jp.client.Issue.Search(jp.filter, sopts)
	s.Stop()
	if err != nil {
		return nil, err
	}

	res := make([]IssueDescription, len(raw))
	for i, r := range raw {
		res[i] = IssueDescription{
			ID:          r.Key,
			Name:        r.Fields.Summary,
			Description: r.Fields.Description,
			Type:        r.Fields.Type.Name,
			Provider:    JIRA_PROVIDER,
		}
	}

	return res, nil
}

func jiraProvider(c *git.Config) (*JiraProvider, error) {
	enable, _ := c.LookupBool("jira.enable")
	if !enable {
		return nil, nil
	}

	username, err := c.LookupString("jira.username")
	if err != nil {
		return nil, err
	}
	token, err := c.LookupString("jira.token")
	if err != nil {
		return nil, err
	}
	domain, err := c.LookupString("jira.domain")
	if err != nil {
		return nil, err
	}
	filter, err := c.LookupString("jira.filter")
	if err != nil {
		return nil, err
	}

	tp := jira.BasicAuthTransport{
		Username: username,
		Password: token,
	}
	jc, err := jira.NewClient(tp.Client(), domain)
	if err != nil {
		return nil, err
	}

	return &JiraProvider{client: jc, filter: filter}, nil
}
