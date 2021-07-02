package integrations

import (
	git "github.com/libgit2/git2go/v31"
)

// Provider interface abstracts cross-platform providers
type Provider interface {
	// Search a list of issue in the provider
	Search() ([]IssueDescription, error)
}

func ProvidersFrom(r *git.Repository) ([]Provider, error) {
	c, err := r.Config()
	if err != nil {
		return nil, err
	}

	var p []Provider

	// Jira
	jp, err := jiraProvider(c)
	if err != nil {
		return nil, err
	}
	if jp != nil {
		p = append(p, *jp)
	}

	// Gitlab
	glp, err := NewGitLabProvider(r)
	if err != nil {
		return nil, err
	}
	if glp != nil {
		p = append(p, *glp)
	}

	return p, nil
}
