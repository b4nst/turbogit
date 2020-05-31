package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
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

func stageNewFile(r *git.Repository) error {
	// Create and stage file
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	f, err := ioutil.TempFile(wd, "*")
	if err != nil {
		return err
	}
	wt, err := r.Worktree()
	if err != nil {
		return err
	}
	_, err = wt.Add(filepath.Base(f.Name()))
	if err != nil {
		return err
	}
	return nil
}

func lastTagFrom(r *git.Repository) (*plumbing.Reference, error) {
	tags, err := r.Tags()
	if err != nil {
		return nil, err
	}

	var tag *plumbing.Reference
	tags.ForEach(func(t *plumbing.Reference) error {
		tag = t
		return nil
	})

	return tag, nil
}
