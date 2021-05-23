package git

import (
	"net/url"
	"regexp"

	git "github.com/libgit2/git2go/v30"
	giturls "github.com/whilp/git-urls"
)

// ParseRemote parses a the remote into an URL.
// If fallback is true and the repository name cannot be found, the first remote found is used
func ParseRemote(r *git.Repository, name string, fallback bool) (*url.URL, error) {
	remote, err := r.Remotes.Lookup(name)
	if err != nil {
		if !fallback {
			return nil, err
		}
		rl, err := r.Remotes.List()
		if err != nil {
			return nil, err
		}
		remote, err = r.Remotes.Lookup(rl[0])
		if err != nil {
			return nil, err
		}
	}
	return giturls.Parse(remote.Url())
}

// Remotes returns the list of repository's remotes
func Remotes(r *git.Repository) ([]git.Remote, error) {
	names, err := r.Remotes.List()
	if err != nil {
		return nil, err
	}
	remotes := make([]git.Remote, len(names))
	for i, rn := range names {
		if remote, err := r.Remotes.Lookup(rn); err == nil {
			remotes[i] = *remote
		}
	}
	return remotes, nil
}

// AnyRemoteMatch return true if any remote url of the repository matches the regexp
func AnyRemoteMatch(r *git.Repository, reg *regexp.Regexp) (bool, error) {
	if reg == nil {
		return false, nil
	}
	remotes, err := Remotes(r)
	if err != nil {
		return false, err
	}
	for _, remote := range remotes {
		if reg.MatchString(remote.Url()) {
			return true, nil
		}
	}
	return false, nil
}
