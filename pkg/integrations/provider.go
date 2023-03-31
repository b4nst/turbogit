package integrations

import (
	git "github.com/libgit2/git2go/v33"
)

// Issuer interface abstracts cross-platform providers
type Issuer interface {
	// Search a list of issue in the provider
	Search() ([]IssueDescription, error)
}

type Commiter interface {
	// Propose commit messages from a diff
	CommitMessages(*git.Diff) ([]string, error)
}

func Issuers(r *git.Repository) (issuers []Issuer, err error) {
	// Jira
	jp, err := NewJiraProvider(r)
	if err != nil {
		return nil, err
	}
	if jp != nil {
		issuers = append(issuers, *jp)
	}

	// Gitlab
	glp, err := NewGitLabProvider(r)
	if err != nil {
		return nil, err
	}
	if glp != nil {
		issuers = append(issuers, *glp)
	}

	return
}

func Commiters(r *git.Repository) (commiters []Commiter, err error) {
	// OpenAI
	oai, err := NewOpenAIProvider(r)
	if err != nil {
		return nil, err
	}
	if oai != nil {
		commiters = append(commiters, oai)
	}

	return
}
