package git

import (
	"net/url"

	git "github.com/libgit2/git2go/v31"
	giturls "github.com/whilp/git-urls"
)

func ParseRemote(r *git.Repository, name string, fallback bool) (*url.URL, error) {
	rawurl := ""
	remote, err := r.Remotes.Lookup(name)
	if err != nil {
		if !fallback {
			return nil, err
		}
		rl, err := r.Remotes.List()
		if err != nil {
			return nil, err
		}
		rawurl = rl[0]
	} else {
		rawurl = remote.Url()
	}
	return giturls.Parse(rawurl)
}
