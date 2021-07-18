# Installation

## macOS

`turbogit` is available on [MacPorts](https://www.macports.org/install.php) or homebrew.

### MacPorts

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

#### Alpine

Install and upgrade:

Download the .apk file from the [releases page](https://github.com/b4nst/turbogit/releases/latest)

```shell
apk add --allow-untrusted ./turbogit_*_linux_amd64.apk
```

## Windows

> Since version 2.0.0 tug is not available as prebuilt binat or package
> for Windows. Please check [#48](https://github.com/b4nst/turbogit/issues/48)
> for more details

`turbogit` is available via [scoop](https://scoop.sh).

```shell
scoop bucket add scoop-bucket https://github.com/b4nst/scoop-bucket.git
scoop install turbogit
```

## Other platforms

Prebuilt binary available from the [release page](https://github.com/b4nst/turbogit/releases/latest).
