# Installation

[![Packaging status](https://repology.org/badge/vertical-allrepos/turbogit.svg)](https://repology.org/project/turbogit/versions)

## macOS

`turbogit` is available on [MacPorts](https://www.macports.org/install.php) and Homebrew.

### Macports (preferred)

install

```shell
sudo port install turbogit
```

upgrade

```shell
sudo port selfupdate && sudo port upgrade turbogit
```

### Homebrew

install

```shell
brew tap b4nst/homebrew-tap
brew install turbogit
```

upgrade

```shell
brew upgrade turbogit
```

## Linux

### NixOS

```
nix-env -i turbogit
```

### Other distributions

Download pre-built binaries from the [latest release page](https://github.com/b4nst/turbogit/releases/latest).

## Windows

> Since git2go refactor, tug is not available as a Windows package anymore.
> Please check [#48](https://github.com/b4nst/turbogit/issues/48).

Please follow [the instructions](/installation#build-from-source) to build it from source.

## Build from source

1. Clone the repo (don't forget the submodule)
2. Run build command

```shell
git clone --recurse-submodules https://github.com/b4nst/turbogit.git
cd turbogit
make build
```
