# Turbogit (tug)

[![build](https://github.com/b4nst/turbogit/workflows/Go/badge.svg)](https://github.com/b4nst/turbogit/actions?query=workflow%3AGo)
[![version](https://img.shields.io/github/v/release/b4nst/turbogit?include_prereleases&label=latest&logo=ferrari)](https://github.com/b4nst/turbogit/releases/latest)
![dependabot](https://api.dependabot.com/badges/status?host=github&repo=b4nst/turbogit)

tug is a cli tool built to help you deal with your day-to-day git work. tug enforces convention (e.g. [The Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)) but tries to keep things simple and invisible for you. tug is your friend.


![logo](assets/tu_logo.png)


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

## Instal

Prebuilt binary available on the [release page](https://github.com/b4nst/turbogit/releases/latest).