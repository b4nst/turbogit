package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

// Hooks

func hookCmd(root string, hook string) (*exec.Cmd, error) {
	script := path.Join(root, ".git", "hooks", hook)
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
		Dir:    root,
		Path:   script,
		Args:   []string{script},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}, nil
}

func noArgHook(root string, hook string) error {
	cmd, err := hookCmd(root, hook)
	if err != nil {
		return err
	}
	if cmd == nil {
		return nil
	}

	fmt.Printf("Running %s hook...\n", hook)
	return cmd.Run()
}

func fileHook(root string, hook string, initial string) (out string, err error) {
	out = initial
	cmd, err := hookCmd(root, hook)
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
	fmt.Printf("Running %s hook...\n", hook)
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

func PreCommitHook(root string) error {
	return noArgHook(root, "pre-commit")
}

func PostCommitHook(root string) error {
	return noArgHook(root, "post-commit")
}

func PrepareCommitMsgHook(root string) (msg string, err error) {
	return fileHook(root, "prepare-commit-msg", "")
}

func CommitMsgHook(root string, in string) (msg string, err error) {
	return fileHook(root, "commit-msg", in)
}
