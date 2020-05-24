<!-- omit in toc -->
# Turbogit (tug)

[![build](https://github.com/b4nst/turbogit/workflows/Go/badge.svg)](https://github.com/b4nst/turbogit/actions?query=workflow%3AGo)
[![version](https://img.shields.io/github/v/release/b4nst/turbogit?include_prereleases&label=latest&logo=ferrari)](https://github.com/b4nst/turbogit/releases/latest)
[![Test Coverage](https://api.codeclimate.com/v1/badges/5173f55b5e67109d3ca5/test_coverage)](https://codeclimate.com/github/b4nst/turbogit/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/5173f55b5e67109d3ca5/maintainability)](https://codeclimate.com/github/b4nst/turbogit/maintainability)
![dependabot](https://api.dependabot.com/badges/status?host=github&repo=b4nst/turbogit)

![logo](assets/tu_logo.png)

tug is a cli tool built to help you deal with your day-to-day git work. tug enforces convention (e.g. [The Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)) but tries to keep things simple and invisible for you. tug is your friend.

![tug example](assets/tug.gif)

- [Usage](#usage)
- [Installation](#installation)
  - [macOS](#macos)
  - [Linux](#linux)
    - [Debian/Ubuntu](#debianubuntu)
    - [Fedora](#fedora)
    - [Centos](#centos)
    - [openSUSE/SUSE](#opensusesuse)
  - [Windows](#windows)
  - [Other platforms](#other-platforms)
- [Shell completion](#shell-completion)
  - [Bash](#bash)
  - [Zsh](#zsh)
  - [Fish](#fish)
- [Contributing](#contributing)

## Usage

```shell
Usage:
  tug [command]

Available Commands:
  branch      Create a new branch.
  commit      Create a new commit.
  config      Read or write config.
  help        Help about any command

Flags:
      --config string   config file (default is $HOME/.config/tug/tug.toml)
  -h, --help            help for tug
```

## Installation

### macOS
`turbogit` is available via Homebrew.

install
```shell
brew tap b4nst/homebrew-tap
brew install turbogit
```

upgrade
```shell
brew upgrade turbogit
```

### Linux
#### Debian/Ubuntu
Install and upgrade:

Download the .deb file from the [releases page](https://github.com/b4nst/turbogit/releases/latest)
```shell
sudo apt install ./turbogit_*_linux_amd64.deb
```
#### Fedora
Install and upgrade:

Download the .rpm file from the [releases page](https://github.com/b4nst/turbogit/releases/latest)
```shell
sudo dnf install ./turbogit_*_linux_amd64.rpm
```
#### Centos
Install and upgrade:

Download the .rpm file from the [releases page](https://github.com/b4nst/turbogit/releases/latest)
```shell
sudo yum localinstall ./turbogit_*_linux_amd64.rpm
```
#### openSUSE/SUSE
Install and upgrade:

Download the .rpm file from the [releases page](https://github.com/b4nst/turbogit/releases/latest)
```shell
sudo zypper in ./turbogit_*_linux_amd64.rpm
```
### Windows
`turbogit` is available via [scoop](https://scoop.sh).
```shell
scoop bucket add scoop-bucket https://github.com/b4nst/scoop-bucket.git
scoop install turbogit
```
### Other platforms
Prebuilt binary available from the [release page](https://github.com/b4nst/turbogit/releases/latest).

## Shell completion

### Bash
```bash
source <(tug completion bash)
```
To load completions for each session, execute once:
*Linux*
```bash
tug completion bash > /etc/bash_completion.d/tug
```
*MacOS*
```bash
tug completion bash > /usr/local/etc/bash_completion.d/tug
```

### Zsh
```zsh
source <(tug completion zsh)
```
To load completions for each session, execute once:
```zsh
tug completion zsh > "${fpath[1]}/_tug"
```

### Fish
```shell
tug completion fish | source
```
To load completions for each session, execute once:
```shell
tug completion fish > ~/.config/fish/completions/tug.fish
```

## Contributing

Thanks for getting involved in turbogit ❤️! Please read the [Contributing guide](CONTRIBUTING.md), it will help you make your awesome contribution.