# Automated Pairing Messages with Git Hooks

Use this `commit-msg` git hook to automatically share credit with your collaborators on a commit, without needing to manually paste strings into your commit message.

Instead, credit them in natural language within the context of your commit message. For example:

```
Add new features. Pairing with @mona.
```

On a system with `gh` and the `gh-pairing-with` extension installed, and in a repository with the `commit-msg` hook enabled, the above commit message will be automatically rewritten to include the appropriate `Co-authored-by:` string.

See examples of supported language in [`scripts/tests/commit_msg_pairs_test.rb`](../scripts/tests/commit_msg_pairs_test.rb). Pull requests welcome!

## What is a Git Hook?

Git Hooks are scripts that are [run automatically when certain actions are taken](https://git-scm.com/docs/githooks) in a git repository. They are not transferred automatically when a repository is cloned, so setting them up requires affirmative action on the user's part.

To be run by git, a git hook must have its executable bit set (e.g. `chmod +x commit-msg`).

You have two options for installing git hooks in your development environment, though they are mutually-exclusive:

1. Adding individual scripts to the `.git/hooks` directory of a given repository.
2. Configuring the `hooksPath` option in your local `.gitconfig` file. This will set a given directory on your system as the overriding path for git hooks in every repository that you interact with. Setting `hooksPath` in your config will cause git to ignore any hooks defined in the repository's `.git/hooks` directory.

Each script must be named to match the name of the action on which occasion it will run.

## Installation

Copy the `commit-msg` script in this directory to the `.git/hooks` directory in the repository where you'd like to use it. Alternately, create a directory (e.g. `~/githooks`), move the `commit-msg` script there, and add the following to your `.gitconfig`:

```
[core]
	hooksPath = ~/githooks
```

### Prerequisites
This script won't work as written without an existing installation of `gh` and the `gh-pairing-with` extension. See [README.md](../README.md#installation) for installation instructions.

