# Usage

## tug

Improve your git workflow.

### Synopsis

Improve your git workflow.

### Options

```
      --config string   config file (default is $HOME/.config/tug/tug.toml)
  -h, --help            help for tug
```

### SEE ALSO

* [tug branch](#tug-branch)	 - Create a new branch
* [tug commit](#tug-commit)	 - Commit staging area
* [tug completion](#tug-completion)	 - Generate completion script
* [tug config](#tug-config)	 - Read or write config
* [tug tag](#tug-tag)	 - Create a tag
* [tug version](#tug-version)	 - Print current version



## tug branch

Create a new branch

### Synopsis

Create a new branch

```
tug branch [feat fix user] [description]
```

### Examples

```

# Create branch feat/my-feature from current branch
$ tug branch feat my feature

# Create branch user/alice/my-branch, given that alice is the current tug/git user
$ tug branch user my branch
	
```

### Options

```
  -h, --help   help for branch
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/tug/tug.toml)
```

### SEE ALSO

* [tug](#tug)	 - Improve your git workflow.



## tug commit

Commit staging area

### Synopsis

Commit staging area

```
tug commit [type] [subject]
```

### Examples

```

# Commit a new feature (feat: a new feature)
$ tug commit feat a new feature

# Commit a fix that brings breaking changes (fix!: API break)
$ tug commit fix -c API break

# Add a scope to the commit (refactor(scope): a scopped refactor)
$ tug commit refactor a scopped refactor -s scope

# Open your editor to edit the commit message
$ tug commit ci -e message
	
```

### Options

```
  -c, --breaking-changes   Commit contains breaking changes
  -e, --edit               Prompt editor to edit your message (add body or/and footer(s))
  -h, --help               help for commit
  -s, --scope string       Add a scope
  -t, --type string        Commit types [build ci chore docs feat fix perf refactor style test]
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/tug/tug.toml)
```

### SEE ALSO

* [tug](#tug)	 - Improve your git workflow.



## tug completion

Generate completion script

### Synopsis

Generate completion script

```
tug completion [bash zsh fish powershell]
```

### Examples

```

Bash:

$ source <(tug completion bash)

# To load completions for each session, execute once:
Linux:
  $ tug completion bash > /etc/bash_completion.d/tug
MacOS:
  $ tug completion bash > /usr/local/etc/bash_completion.d/tug

Zsh:

$ source <(tug completion zsh)

# To load completions for each session, execute once:
$ tug completion zsh > "${fpath[1]}/_tug"

Fish:

$ tug completion fish | source

# To load completions for each session, execute once:
$ tug completion fish > ~/.config/fish/completions/tug.fish

```

### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/tug/tug.toml)
```

### SEE ALSO

* [tug](#tug)	 - Improve your git workflow.



## tug config

Read or write config

### Synopsis

If [value] is provided, sets [key] to [value] in config, otherwise print current value for [key]

```
tug config [key] [value]
```

### Examples

```

# Set config user.name to alice. If the config does not exist, it will be created.
$ tug config user.name alice

# Get the current value for user.name
$ tug config user.name

# Delete the entrie user.name
$ tug config -d user.name

```

### Options

```
  -d, --delete   Delete config
  -h, --help     help for config
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/tug/tug.toml)
```

### SEE ALSO

* [tug](#tug)	 - Improve your git workflow.



## tug tag

Create a tag

### Synopsis

Create a semver tag, based on the commit history since last one

```
tug tag
```

### Examples

```

# Given that the last release tag was v1.0.0, some feature were committed but no breaking changes.
# The following command will create the tag v1.1.0
$ tug tag	

```

### Options

```
  -d, --dry-run   Do not tag.
  -h, --help      help for tag
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/tug/tug.toml)
```

### SEE ALSO

* [tug](#tug)	 - Improve your git workflow.



## tug version

Print current version

### Synopsis

Print current version

```
tug version
```

### Options

```
  -h, --help   help for version
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.config/tug/tug.toml)
```

### SEE ALSO

* [tug](#tug)	 - Improve your git workflow.


