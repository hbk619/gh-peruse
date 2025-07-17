# Git Peruse

For using git in a more screen reader friendly way and for those who like to go one at a time.

You can use the binary as a [Github cli extension](#github-cli-extension) (preferred way) or as a [git command](#without-installing-github-cli),
the latter will require a Personal Access token to be created.

## Pre-requisites (Linux)
If you are using Linux and want to use the "copy comment to clipboard" feature you will need
[xclip](https://github.com/astrand/xclip) or [xsel](https://github.com/kfish/xsel) installed. These are generally available from your package manager.

## Github CLI extension
Install the [Github CLI](https://cli.github.com/)

Run `gh auth login` and approve the app in the browser. 

### Install the extension

`gh extension install https://github.com/hbk619/gh-peruse`

### Usage

To view the current PR for your branch:

`gh peruse pr`

To view a specific PR:

`gh peruse pr <pr number>`

For full up-to-date flags:

`gh peruse pr -h`

## Without installing Github CLI
`gh-peruse` uses the Github CLI library and follows the authentication mechanism and configuration options it offers:

```markdown
GitHub API requests will be authenticated using the same mechanism as gh, i.e. using the values of GH_TOKEN and GH_HOST environment variables and falling back to the user's stored OAuth token
```

You can use `gh-peruse` without installing the Github CLI by setting a [personal access token](#personal-access-token) or the [install the Github CLI](#github-cli-extension) to use OAuth.

### Personal Access token
Create a [classic personal access token](https://github.com/settings/tokens/new) with the repo scope
and put it in an environment variable called `GH_TOKEN` or `GITHUB_TOKEN`.

If you are using an enterprise server use `GH_ENTERPRISE_TOKEN` or `GITHUB_ENTERPRISE_TOKEN`

### Changing the host
If the host cannot be inferred from the context of a repo (shouldn't happen often), set an environment variable called `GH_HOST`

## Gitlab
If you're using Gitlab, you will need a [personal access token](#gitlab-personal-access-token) and [environment variable](#set-environment-variables), then [install for use with git](#install-the-binary-for-use-with-git)

### Gitlab Personal Access Token

Create a [personal access token](https://gitlab.com/-/user_settings/personal_access_tokens?page=1&state=active&sort=expires_asc) with scopes `api`

### Set environment Variables
You need two variables:
- GITLAB_TOKEN - your gitlab personal access token
- PERUSE_TYPE - to force Peruse to use Gitlab APIs

How to set varies per operating system and between versions of operating systems, sorry!

Linux/Mac can be set temporarily with `export GITLAB_TOKEN=<my-token-here>`

### Install the binary for use with Git

#### MacOS and Linux

Download gh-peruse-darwin-amd64 for Mac or gh-peruse-linux-amd64 for linux from
[latest release](https://github.com/hbk619/gh-peruse/releases/latest) 

Copy the binary to somewhere on your path (e.g. `/usr/local/bin`) called `git-peruse`

#### Windows

Download gh-peruse-windows-amd64.exe for Windows from
[latest release](https://github.com/hbk619/gh-peruse/releases/latest) 

Copy the binary to somewhere on your path and name it `git-peruse.exe`

### Usage

To view the current PR for your branch:

`git peruse pr`

To view a specific PR:

`git peruse pr <pr number>`

For full up-to-date flags:

`git peruse pr -h`

### Notifications for new comments

If you'd like to receive a system notification when you have new comments on PRs you own you can add the following to a job that runs on a schedule e.g. via [cron](https://en.wikipedia.org/wiki/Cron) or [Windows Task Scheduler](https://en.wikipedia.org/wiki/Windows_Task_Scheduler)

#### Using Github CLI

`gh peruse pr check -n`

#### Using Git

`git peruse pr check -n`

## Developing

Install [Go](https://go.dev/doc/install)

### Building

`go build -o gh-peruse main.go`

This gives a binary in the dir you run the command from called 'gh-peruse' which you can run with:

`./gh-peruse pr 1`

### Tests

Run with `go test ./...`