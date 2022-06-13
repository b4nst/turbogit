# Installation

[![Packaging status](https://repology.org/badge/vertical-allrepos/turbogit.svg)](https://repology.org/project/turbogit/versions)

> turbogit needs [libgit2](https://github.com/libgit2/libgit2) >= 1.1.0 on your system.
> Some package manager will handle its installation automatically

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

### Debian/Ubuntu

Install and upgrade:

Download the .deb file from the [releases page](https://github.com/b4nst/turbogit/releases/latest)

```shell
sudo apt install ./turbogit_*_linux_amd64.deb
```

### Fedora

Install and upgrade:

Download the .rpm file from the [releases page](https://github.com/b4nst/turbogit/releases/latest)

```shell
sudo dnf install ./turbogit_*_linux_amd64.rpm
```

### Centos

Install and upgrade:

Download the .rpm file from the [releases page](https://github.com/b4nst/turbogit/releases/latest)

```shell
sudo yum localinstall ./turbogit_*_linux_amd64.rpm
```

### openSUSE/SUSE

Install and upgrade:

Download the .rpm file from the [releases page](https://github.com/b4nst/turbogit/releases/latest)

```shell
sudo zypper in ./turbogit_*_linux_amd64.rpm
```

### Alpine

Install and upgrade:

Download the .apk file from the [releases page](https://github.com/b4nst/turbogit/releases/latest)

```shell
apk add --allow-untrusted ./turbogit_*_linux_amd64.apk
```

## Windows

> Since git2go refactor, tug is not available as a Windows package anymore.
> Please check [#48](https://github.com/b4nst/turbogit/issues/48).

`turbogit` is available via [scoop](https://scoop.sh).

```shell
scoop bucket add scoop-bucket https://github.com/b4nst/scoop-bucket.git
scoop install turbogit
```

## Other platforms

Prebuilt binary available from the [release page](https://github.com/b4nst/turbogit/releases/latest).
