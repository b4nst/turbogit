package integrations

import (
	"fmt"
	"strconv"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/briandowns/spinner"
	"github.com/go-git/go-git/v5/plumbing/format/config"
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

func jiraProvider(c *config.Config) (*JiraProvider, error) {
	s := c.Section("jira")
	enable, err := strconv.ParseBool(s.Option("enable"))
	if err != nil {
		return nil, fmt.Errorf("bad config format: %w", err)
	}
	if !enable {
		return nil, nil
	}

	username := s.Option("username")
	token := s.Option("token")
	domain := s.Option("domain")
	filter := s.Option("filter")
	fmt.Println("Jira", enable, username, token, domain, filter)

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
