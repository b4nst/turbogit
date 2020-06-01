package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/b4nst/turbogit/cmd"
	"github.com/spf13/cobra/doc"
)

const (
	docDir    = "dist/doc"
	assetsDir = "assets"
)

func main() {
	if err := genDocsify(); err != nil {
		fatal(err)
	}
}

func genDocsify() error {
	outdir := path.Join(docDir, "docsify")
	sidebar := []string{}

	filePrepender := func(filename string) string {
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		sidebar = append(sidebar, fmt.Sprintf("* [%s](/%s)", strings.Replace(base, "_", " ", -1), name))
		return ""
	}
	identity := func(s string) string { return s }
	if err := os.MkdirAll(outdir, 0755); err != nil {
		return err
	}
	err := doc.GenMarkdownTreeCustom(cmd.RootCmd, outdir, filePrepender, identity)
	if err != nil {
		return err
	}
	// Reorder sidebar
	x, sidebar := sidebar[len(sidebar)-1], sidebar[:len(sidebar)-1]
	sidebar = append([]string{x}, sidebar...)
	// Write sidebar
	sbF, err := os.Create(filepath.Join(outdir, "_sidebar.md"))
	if err != nil {
		return err
	}
	defer sbF.Close()
	sbF.WriteString(strings.Join(sidebar, "\n"))
	// Touch .nojekyll
	njF, err := os.Create(filepath.Join(outdir, ".nojekyll"))
	if err != nil {
		return err
	}
	defer njF.Close()

	// Copy assets
	// index.html
	_, err = copy(filepath.Join(assetsDir, "doc", "index.html"), filepath.Join(outdir, "index.html"))
	if err != nil {
		return err
	}
	// _coverpage.md
	_, err = copy(filepath.Join(assetsDir, "doc", "_coverpage.md"), filepath.Join(outdir, "_coverpage.md"))
	if err != nil {
		return err
	}
	// icon.png
	if err := os.MkdirAll(filepath.Join(outdir, "_media"), 0755); err != nil {
		return err
	}
	_, err = copy(filepath.Join(assetsDir, "tu_logo.png"), filepath.Join(outdir, "_media", "icon.png"))
	if err != nil {
		return err
	}

	return nil
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func fatal(msg interface{}) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
