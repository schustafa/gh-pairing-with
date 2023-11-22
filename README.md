# `gh-pairing-with`

A simple [`gh`](https://github.com/cli/cli) extension for sharing credit.

## Background

Did you know that GitHub can show multiple authors in commit statistics? For example, see the [early history of this repository](https://github.com/schustafa/gh-pairing-with/commits/main?after=31df6e42fb6e7801b05285736da6d0b5a5431595+1&branch=main):

<img width="811" alt="Early commits to this repository, showing contributions from both @schustafa and @stephanieg0" src="https://github.com/schustafa/gh-pairing-with/assets/126731/41e1060b-7ce7-45ba-b26d-452bae06282d">

For GitHub to know what users besides the author to credit, you must [include a special `Co-authored-by:` string in your commit message](https://docs.github.com/en/pull-requests/committing-changes-to-your-project/creating-and-editing-commits/creating-a-commit-with-multiple-authors). Unfortunately, the format is finicky (make sure your hyphens are in the right place!) and you need to correctly include the name and email address associated with your coauthor's GitHub account.

I lost count of how many times I fat-fingered this string, committed, pushed, and then realized that I'd gotten something wrong. This extension exists to make getting it right slightly easier.

## Installation

Install the CLI

```bash
brew install gh
```

Install the extension

```bash
gh extension install schustafa/gh-pairing-with
```

## Commands

Run `gh pairing-with <github-login>...`.

For example, if you're pairing with [Miss Monalisa Octocat](https://github.com/mona), you'll run:

```bash
> gh pairing-with mona
Co-authored-by: Monalisa Octocat <92997159+mona@users.noreply.github.com>
```

Paste the string returned into your commit message to share credit with your pairing partner!

If you're on a Mac, pipe the output to `pbcopy` to get the `Co-authored-by` string automatically added to your pasteboard, ready to paste!

```bash
gh pairing-with mona | pbcopy
```

`gh-pairing-with` supports fetching user information for multiple users at once too! So if you're, uh, "pairing" with a handful of people, you can pass multiple usernames at once. For example:

```bash
> gh pairing-with mona schustafa
Co-authored-by: Monalisa Octocat <92997159+mona@users.noreply.github.com>
Co-authored-by: AJ Schuster <126731+schustafa@users.noreply.github.com>
```

## Troubleshooting

### Scopes

If you receive a message like the following:

```
GraphQL: Your token has not been granted the required scopes to execute this query. The 'email' field requires one of the following scopes: ['user:email', 'read:user'], but your token has only been granted the: [...]
```

You can add those scopes to your `gh` token by running the following:

```bash
gh auth refresh --scopes user:email,read:user
```

## Resources

- [Setting your commit email address](https://docs.github.com/en/account-and-profile/setting-up-and-managing-your-personal-account-on-github/managing-email-preferences/setting-your-commit-email-address)
- [Your GitHub Email Settings](https://github.com/settings/emails)
