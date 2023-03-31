# Integration

Turbogit can plug itself with with some of your favourite tools whenever it makes sense.
For instance including a ticket id in the branch name.
If you think something is missing in your workflow with turbogit, do not hesitate to raise an issue on [b4nst/turbogit](https://github.com/b4nst/turbogit/issues).

## OpenAI integration

The OpenAI integration enables you to fill commit messages automatically based on the staged diff.
It currently uses `gpt-3.5-turbo` model.

### Configuration

In order to enable OpenAI integration, you will need to set some configuration keys.

| key                            | type     | description                                                                                                                         | recommended location |
| ---                            | ---      | ---                                                                                                                                 | ---                  |
| **openai.enabled**             | `bool`   | Enable GitLab integration                                                                                                           | global               |
| **openai.token**               | `string` | OpenAI [API key](https://platform.openai.com/account/api-keys)                                                                      | global               |


Set a global key:

```shell
git config --global <key> <value>
```

Set a local key

```shell
git config <key> <value>
```

### Usage

First set the proper configuration to activate OpenAI integration.
Then you just have to

```shell
tug commit --fill
```

This will ask OpenAI for a commit message based on the current staged diff.
You can refine it either overwritting by passing other args (type, summary, scope, etc) or
spawn a editor with `-e` option.


## GitHub integration

_This is a work in progress, please check [#60](https://github.com/b4nst/turbogit/issues/60) for further details._

## GitLab integration

The Gitlab integration enables you to create branches automatically from GitLab issues.

### Configuration

In order to enable GitLab integration, you will need to set some configuration keys.

| key                            | type     | description                                                                                                                         | recommended location |
| ---                            | ---      | ---                                                                                                                                 | ---                  |
| **gitlab.enabled**             | `bool`   | Enable GitLab integration                                                                                                           | local                |
| **gitlab.token**               | `string` | GitLab [personal access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#create-a-personal-access-token). | global               |
| **gitlab.protocol** (optional) | `string` | Override GitLab API protocol (default https)                                                                                        | -                    |


Set a global key:

```shell
git config --global <key> <value>
```

Set a local key

```shell
git config <key> <value>
```

### Usage

First set the proper configuration to activate GitLab integration.
Then you just have to

```shell
tug new
```

It will prompt you a list of issues with a fuzzy finder at your disposal to refine your selection.
Select your issue and turbogit will take care of creating and checkout the branch for you.

## Jira integration

The Jira integration enables you to create branches automatically from Jira issues.

### Configuration

In order to enable Jira integration, you will need to set those configuration keys:

| key               | type     | description                                                                                         | recommended location                                        |
| ---               | ---      | ---                                                                                                 | ---                                                         |
| **jira.enabled**   | `bool`   | Enable jira integration                                                                             | local                                                       |
| **jira.token**    | `string` | Jira personal token. Create one [here](https://id.atlassian.com/manage-profile/security/api-tokens) | global                                                      |
| **jira.username** | `string` | Your Jira username (email)                                                                          | global                                                      |
| **jira.domain**   | `string` | Your Jira domain, including protocol (e.g. https://company.atlassian.net)                           | global                                                      |
| **jira.filter**   | `string` | JQL filter to gather issues                                                                         | global: a wide filter, local: override with narrower filter |

Set a global key:

```shell
git config --global <key> <value>
```

Set a local key

```shell
git config <key> <value>
```

### Usage

First set the proper configuration to activate Jira integration.
Then you just have to

```shell
tug new
```

It will prompt you a list of issues matching `jira.filter` with a fuzzy finder at your disposal to refine your selection.
Select your issue and turbogit will take care of creating and checkout the branch for you.
