package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

func setUpRepo() (r *git.Repository, teardown func(), err error) {
	dir, err := ioutil.TempDir("", "turbogit-test")
	if err != nil {
		return nil, func() {}, err
	}
	teardown = func() {
		os.RemoveAll(dir)
	}

	if err := os.Chdir(dir); err != nil {
		return nil, teardown, err
	}

	d, err := os.Getwd()
	if err != nil {
		return nil, teardown, err
	}

	r, err = git.PlainInit(d, false)
	if err != nil {
		return nil, teardown, err
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, teardown, err
	}

	// Create a HEAD
	wd, err := os.Getwd()
	if err != nil {
		return nil, teardown, err
	}

	filename := filepath.Join(wd, "example-git-file")
	err = ioutil.WriteFile(filename, []byte("hello world!"), 0644)
	if err != nil {
		return nil, teardown, err
	}
	_, err = w.Add("example-git-file")
	if err != nil {
		return nil, teardown, err
	}

	cfg, err := r.ConfigScoped(config.SystemScope)
	if err != nil {
		return nil, teardown, err
	}
	cfg.Author.Name = "John Doe"
	cfg.Author.Email = "john@doe.org"
	r.SetConfig(cfg)

	_, err = w.Commit("example go-git commit", &git.CommitOptions{})
	if err != nil {
		return nil, teardown, err
	}

	return r, teardown, nil
}
