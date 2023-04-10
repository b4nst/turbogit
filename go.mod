module github.com/b4nst/turbogit

go 1.14

require (
	github.com/AlecAivazis/survey/v2 v2.3.6
	github.com/andygrunwald/go-jira v1.16.0
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de
	github.com/blang/semver/v4 v4.0.0
	github.com/briandowns/spinner v1.23.0
	github.com/fatih/color v1.15.0 // indirect
	github.com/gdamore/tcell/v2 v2.6.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-hclog v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hpcloud/golor v0.0.0-20150914221010-dc1b58c471a0
	github.com/imdario/mergo v0.3.15
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/ktr0731/go-fuzzyfinder v0.7.0
	github.com/libgit2/git2go/v33 v33.0.9
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/sashabaranov/go-openai v1.7.0
	github.com/spf13/cobra v1.6.1
	github.com/stretchr/testify v1.8.2
	github.com/whilp/git-urls v1.0.0
	github.com/xanzy/go-gitlab v0.81.0
	golang.org/x/crypto v0.7.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/libgit2/git2go/v33 => ./git2go
