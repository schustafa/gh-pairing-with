# `gh-pairing-with`

A simple [`gh`](https://github.com/cli/cli) extension for sharing credit.

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

Run `gh pairing-with <github-login>`.

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

## Resources

- [Setting your commit email address](https://docs.github.com/en/account-and-profile/setting-up-and-managing-your-personal-account-on-github/managing-email-preferences/setting-your-commit-email-address)
- [Your GitHub Email Settings](https://github.com/settings/emails)
