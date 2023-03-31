module github.com/b4nst/turbogit

go 1.14

require (
	github.com/AlecAivazis/survey/v2 v2.3.6
	github.com/andygrunwald/go-jira v1.16.0
	github.com/araddon/dateparse v0.0.0-20201001162425-8aadafed4dc4
	github.com/blang/semver/v4 v4.0.0
	github.com/briandowns/spinner v1.19.0
	github.com/fatih/color v1.13.0 // indirect
	github.com/hashicorp/go-hclog v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hpcloud/golor v0.0.0-20150914221010-dc1b58c471a0
	github.com/imdario/mergo v0.3.12
	github.com/kr/text v0.2.0 // indirect
	github.com/ktr0731/go-fuzzyfinder v0.7.0
	github.com/libgit2/git2go/v33 v33.0.9
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/sashabaranov/go-openai v1.5.7
	github.com/spf13/cobra v1.4.1-0.20220414043027-bf6cb5804d7a
	github.com/stretchr/testify v1.8.1
	github.com/whilp/git-urls v1.0.0
	github.com/xanzy/go-gitlab v0.74.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/libgit2/git2go/v33 => ./git2go
