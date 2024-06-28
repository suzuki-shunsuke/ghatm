# ghatm

Set [timeout-minutes](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#jobsjob_idtimeout-minutes) to all GitHub Actions jobs.
`ghatm` finds GitHub Actions workflows and add `timeout-minutes` to jobs which don't have the setting.

`ghatm` edits workflow files but keeps YAML comments, indents, empty lines, and so on.

## Motivation

- https://exercism.org/docs/building/github/gha-best-practices#h-set-timeouts-for-workflows
- [job_timeout_minutes_is_required | suzuki-shunsuke/ghalint](https://github.com/suzuki-shunsuke/ghalint/blob/main/docs/policies/012.md)
- [job_timeout_minutes_is_required | lintnet-modules/ghalint](https://github.com/lintnet-modules/ghalint/tree/main/workflow/job_timeout_minutes_is_required)

`timeout-minutes` should be set properly, but if you have a lot of workflows which don't set `timeout-minutes` it's so bothersome to fix all of them by hand.
`ghatm` sets `timeout-minutes` automatically.

## How to install

`ghatm` is a single binary written in Go.
So you only need to put the executable binary into `$PATH`.

1. [Homebrew](https://brew.sh/)

```sh
brew install suzuki-shunsuke/ghatm/ghatm
```

2. [Scoop](https://scoop.sh/)

```sh
scoop bucket add lintnet https://github.com/suzuki-shunsuke/scoop-bucket
scoop install ghatm
```

3. [aqua](https://aquaproj.github.io/)

```sh
aqua g -i suzuki-shunsuke/ghatm
```

4. Download a prebuilt binary from [GitHub Releases](https://github.com/lintnet/lintnet/releases) and install it into `$PATH`

5. Go

```sh
go install github.com/suzuki-shunsuke/ghatm@latest
```

## How to use

Please run `ghatm set` at the repository root directory.

```sh
ghatm set
```

then `ghatm` checks GitHub Actions workflows `^\.github/workflows/.*\.ya?ml$` and sets `timeout-minutes: 30` to jobs which don't have `timeout-minutes`.
Jobs which have `timeout-minutes` aren't changed.
You can specify the value of `timeout-minutes` with `-t` option.

```sh
ghatm set -t 60
```

You can specify workflow files by positional arguments.

```sh
ghatm set .github/workflows/test.yaml
```

## Tips: Fix workflows by CI

Using `ghatm` in CI, you can fix workflows automatically.
When workflow files are added or changed in a pull request, you can run `ghatm` and commit and push changes to a feature branch.

## LICENSE

[MIT](LICENSE)
