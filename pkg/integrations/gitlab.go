package integrations

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5/config"
	fconfig "github.com/go-git/go-git/v5/plumbing/format/config"
	gurl "github.com/whilp/git-urls"
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

func NewGitLabProvider(c *config.Config) (*GitLabProvider, error) {
	s := c.Raw.Section("gitlab")
	enabled, err := strconv.ParseBool(s.Option("enabled"))
	if err == nil && !enabled {
		return nil, nil
	}
	remote, err := gurl.Parse(c.Remotes["origin"].URLs[0]) // TODO deal with no origin repo
	if err != nil {
		return nil, err
	}
	if !isGitLabRemote(remote, s) {
		if enabled {
			return nil, fmt.Errorf("GitLab provider is enabled but %s is not a known gitlab host. Please add it to gitlab.hosts config or disable GitLab provider for this repository", remote.Hostname())
		}
		return nil, nil
	}
	token := s.Option("token")
	protocol := s.Option("protocol")
	if protocol == "" {
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

func isGitLabRemote(remote *url.URL, s *fconfig.Section) bool {
	hosts := append([]string{GITLAB_CLOUD_HOST}, strings.Split(s.Option("hosts"), ",")...)
	remoteHost := remote.Hostname()
	for _, h := range hosts {
		if h == remoteHost {
			return true
		}
	}
	return false
}
