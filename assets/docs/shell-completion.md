# Shell completion

## Fish

```shell
tug completion fish | source
```

To load completions for each session, execute once:

```shell
tug completion fish > ~/.config/fish/completions/tug.fish
```

## Zsh

```zsh
source <(tug completion zsh)
```

To load completions for each session, execute once:

```zsh
tug completion zsh > "${fpath[1]}/_tug"
```

## Bash

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
