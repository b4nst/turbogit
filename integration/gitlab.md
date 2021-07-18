# GitLab integration

The Gitlab integration enables you to create branches automatically from GitLab issues.

## Configuration

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

## Usage

First set the proper configuration to activate GitLab integration.
Then you just have to

```shell
tug branch
```

It will prompt you a list of issues with a fuzzy finder at your disposal to refine your selection.
Select your issue and turbogit will take care of creating and checkout the branch for you.
