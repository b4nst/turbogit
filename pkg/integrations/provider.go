package integrations

import (
	"fmt"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

// Provider interface abstracts cross-platform providers
type Provider interface {
	// Search a list of issue in the provider
	Search() ([]IssueDescription, error)
}

func ProvidersFrom(r *git.Repository) ([]Provider, error) {
	// c, err := config.LoadConfig(config.GlobalScope)
	c, err := r.ConfigScoped(config.GlobalScope)
	if err != nil {
		return nil, err
	}

	fmt.Println("Jira section", c.Raw.Section("jira").Options.GoString())
	var p []Provider

	// Jira
	jp, err := jiraProvider(c.Raw)
	if err != nil {
		return nil, err
	}
	if jp != nil {
		p = append(p, *jp)
	}

	// Gitlab
	glp, err := NewGitLabProvider(c)
	if err != nil {
		return nil, err
	}
	if glp != nil {
		p = append(p, *glp)
	}

	return p, nil
}
