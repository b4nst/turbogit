package integrations

import (
	"fmt"
	"net/url"
	"strings"

	tugit "github.com/b4nst/turbogit/pkg/git"
	git "github.com/libgit2/git2go/v31"
	"github.com/xanzy/go-gitlab"
)

const (
	// GitLab provider's name
	GITLAB_PROVIDER = "Gitlab"
	// GitLab cloud host
	GITLAB_CLOUD_HOST = "gitlab.com"
	// GitLab default protocol
	GITLAB_DEFAULT_PROTOCOL = "https"
)

type GitLabProvider struct {
	project string
	client  *gitlab.Client
}

// Search return a list of issues of a GitLab project
func (glp GitLabProvider) Search() ([]IssueDescription, error) {
	scope := "assigned_to_me"
	issues, _, err := glp.client.Issues.ListProjectIssues(glp.project, &gitlab.ListProjectIssuesOptions{Scope: &scope})
	if err != nil {
		return nil, err
	}

	res := make([]IssueDescription, len(issues))
	for i, r := range issues {
		res[i] = IssueDescription{
			ID:          fmt.Sprint(r.IID),
			Name:        r.Title,
			Description: r.Description,
			// TODO provide type
			Provider: GITLAB_PROVIDER,
		}
	}
	return res, nil
}

func NewGitLabProvider(r *git.Repository) (*GitLabProvider, error) {
	c, err := r.Config()
	if err != nil {
		return nil, err
	}
	enabled, err := c.LookupBool("gitlab.enabled")
	if err == nil && !enabled {
		return nil, nil
	}
	remote, err := tugit.ParseRemote(r, "origin", true)
	if err != nil {
		return nil, err
	}
	if !isGitLabRemote(remote, c) {
		if enabled {
			return nil, fmt.Errorf("GitLab provider is enabled but %s is not a known gitlab host. Please add it to gitlab.hosts config or disable GitLab provider for this repository", remote.Hostname())
		}
		return nil, nil
	}
	token, err := c.LookupString("gitlab.token")
	if err != nil {
		return nil, err
	}
	protocol, err := c.LookupString("gitlab.protocol")
	if err != nil {
		protocol = GITLAB_DEFAULT_PROTOCOL
	}
	baseUrl := protocol + "://" + remote.Host
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseUrl))
	if err != nil {
		return nil, err
	}
	fmt.Println(remote.Path)
	project := strings.TrimSuffix(remote.Path, ".git")

	return &GitLabProvider{client: client, project: project}, nil
}

func isGitLabRemote(remote *url.URL, c *git.Config) bool {
	hosts := []string{GITLAB_CLOUD_HOST}
	if rhosts, err := c.LookupString("gitlab.hosts"); err == nil {
		hosts = append(hosts, strings.Split(rhosts, ",")...)
	}
	remoteHost := remote.Hostname()
	for _, h := range hosts {
		if h == remoteHost {
			return true
		}
	}
	return false
}
