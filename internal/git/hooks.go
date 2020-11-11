package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/b4nst/turbogit/internal/context"
)

// Hooks

func hookCmd(ctx *context.Context, hook string) (*exec.Cmd, error) {
	wt, err := ctx.Repo.Worktree()
	if err != nil {
		return nil, err
	}
	script := path.Join(wt.Filesystem.Root(), ".git", "hooks", hook)
	info, err := os.Stat(script)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("Hook .git/hooks/%s is a directory, it should be an executable file.", hook)
	}
	return &exec.Cmd{
		Dir:    wt.Filesystem.Root(),
		Path:   script,
		Args:   []string{script},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}, nil
}

func noArgHook(ctx *context.Context, hook string) error {
	cmd, err := hookCmd(ctx, hook)
	if err != nil {
		return err
	}
	if cmd == nil {
		return nil
	}

	return cmd.Run()
}

func fileHook(ctx *context.Context, hook string, initial string) (out string, err error) {
	out = initial
	cmd, err := hookCmd(ctx, hook)
	if cmd == nil {
		return initial, nil
	}

	file, err := ioutil.TempFile("", "file-hook-")
	if err != nil {
		return
	}
	defer file.Close()
	_, err = file.Write([]byte(initial))
	if err != nil {
		return
	}
	file.Close()

	cmd.Args = append(cmd.Args, file.Name())
	err = cmd.Run()
	if err != nil {
		return
	}

	file, err = os.Open(file.Name())
	defer file.Close()
	if err != nil {
		return
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	out = string(content)
	return
}

func PreCommitHook(ctx *context.Context) error {
	fmt.Println("Running pre-commit hook...")
	return noArgHook(ctx, "pre-commit")
}

func PostCommitHook(ctx *context.Context) error {
	fmt.Println("Running post-commit hook...")
	return noArgHook(ctx, "post-commit")
}

func PrepareCommitMsgHook(ctx *context.Context) (msg string, err error) {
	fmt.Println("Running prepare-commit-msg hook...")
	return fileHook(ctx, "prepare-commit-msg", "")
}

func CommitMsgHook(ctx *context.Context, in string) (msg string, err error) {
	fmt.Println("Running commit-msg hook...")
	return fileHook(ctx, "commit-msg", in)
}
