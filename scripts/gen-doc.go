package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/b4nst/turbogit/cmd"
	"github.com/spf13/cobra/doc"
	"gopkg.in/yaml.v3"
)

const (
	Workdir       = "dist/doc"
	AssetsSrcDir  = "assets"
	titleTemplate = `---
title: %s
---
`
)

var (
	DocsDir    = path.Join(Workdir, "docs")
	IncludeDir = path.Join(DocsDir, "_include")
)

func main() {
	log.Println("Ensure output dir exists and is empty...")
	checkErr(os.RemoveAll(Workdir))
	checkErr(os.MkdirAll(Workdir, 0700))

	log.Println("Prepare doctave structure...")
	checkErr(os.MkdirAll(IncludeDir, 0700))
	doctave := Doctave{
		Title:  "Turbogit",
		Colors: Colors{Main: "#e96900"},
		Logo:   "assets/tu_logo.png",
	}

	log.Println("Copy assets...")
	_, err := Copy(path.Join(IncludeDir, doctave.Logo), path.Join(AssetsSrcDir, "tu_logo.png"), "")
	_, err = Copy(path.Join(IncludeDir, "assets/tug.gif"), path.Join(AssetsSrcDir, "tug.gif"), "")
	checkErr(err)

	log.Println("Copy static documentation")
	_, err = Copy(path.Join(DocsDir, "contributing.md"), "CONTRIBUTING.md", "Contributing")
	_, err = Copy(path.Join(DocsDir, "code-of-conduct.md"), "CODE_OF_CONDUCT.md", "Code of conduct")
	_, err = Copy(path.Join(DocsDir, "README.md"), "README.md", "")
	_, err = Copy(path.Join(DocsDir, "installation.md"), "assets/docs/installation.md", "Installation")
	_, err = Copy(path.Join(DocsDir, "usage.md"), "assets/docs/usage.md", "Usage")
	_, err = Copy(path.Join(DocsDir, "integration.md"), "assets/docs/integration.md", "Integration")
	_, err = Copy(path.Join(DocsDir, "shell-completion.md"), "assets/docs/shell-completion.md", "Shell completion")
	checkErr(err)

	cmdDir := path.Join(DocsDir, "Commands")
	log.Println("Ensure commands dir exists...")
	checkErr(os.MkdirAll(cmdDir, 0700))
	log.Println("Generate commands documentation...")
	filePrepender := func(filename string) string {
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		return fmt.Sprintf(titleTemplate, strings.ReplaceAll(base, "_", " "))
	}
	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		return "/Commands/" + strings.ToLower(base)
	}
	checkErr(doc.GenMarkdownTreeCustom(cmd.RootCmd, cmdDir, filePrepender, linkHandler))

	log.Println("Generate nav bar...")
	doctave.Navigation = []Nav{
		{
			Path: "docs/installation.md",
		},
		{
			Path: "docs/usage.md",
		},
		{
			Path:     strings.TrimPrefix(cmdDir, Workdir+"/"),
			Children: "*",
		},
		{
			Path: "docs/integration.md",
		},
		{
			Path: "docs/shell-completion.md",
		},
		{
			Path: "docs/contributing.md",
		},
		{
			Path: "docs/code-of-conduct.md",
		},
	}

	log.Println("Marshal doctave configuration...")
	d, err := yaml.Marshal(doctave)
	checkErr(err)
	checkErr(os.WriteFile(path.Join(Workdir, "doctave.yaml"), d, 0700))

	log.Println("Add GitHub page files...")
	nj, err := os.Create(path.Join(IncludeDir, ".nojekyll"))
	checkErr(err)
	defer nj.Close()
	_, err = Copy(path.Join(IncludeDir, "favicon.ico"), path.Join(AssetsSrcDir, "tu_logo.ico"), "")

	log.Println("Done.")
	log.Println("Use 'doctave serve' or 'doctave build'")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Copy(dst, src, rename string) (int64, error) {
	stat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !stat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	if err := os.MkdirAll(path.Dir(dst), 0700); err != nil {
		return 0, err
	}

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	if rename != "" {
		fmt.Fprintf(destination, titleTemplate, rename)
	}
	defer destination.Close()
	return io.Copy(destination, source)
}

type Doctave struct {
	Title      string `yaml:"title"`
	Port       uint   `yaml:"port,omitempty"`
	BasePath   string `yaml:"base_path,omitempty"`
	DocsDir    string `yaml:"docs_dir,omitempty"`
	Logo       string `yaml:"logo,omitempty"`
	Colors     Colors `yaml:"colors,omitempty"`
	Navigation []Nav  `yaml:"navigation,omitempty"`
}

type Nav struct {
	Path     string `yaml:"path"`
	Children string `yaml:"children,omitempty"`
}

type Colors struct {
	Main string `yaml:"main,omitempty"`
}
