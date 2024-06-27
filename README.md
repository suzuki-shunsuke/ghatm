# gha-set-timeout-minutes

Set [timeout-minutes](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#jobsjob_idtimeout-minutes) to all GitHub Actions jobs

## Motivation

- https://exercism.org/docs/building/github/gha-best-practices#h-set-timeouts-for-workflows
- [job_timeout_minutes_is_required | suzuki-shunsuke/ghalint](https://github.com/suzuki-shunsuke/ghalint/blob/main/docs/policies/012.md)
- [job_timeout_minutes_is_required | lintnet-modules/ghalint](https://github.com/lintnet-modules/ghalint/tree/main/workflow/job_timeout_minutes_is_required)

`timeout-minutes` should be set properly, but if you have a lot of workflows which don't set `timeout-minutes` it's so bothersome to fix all of them by hand.
`gha-set-timeout-minutes` sets `timeout-minutes` automatically.

## Install

`gha-set-timeout-minutes` is a single binary written in Go.
So you only need to put the executable binary into `$PATH`.

```sh
go install github.com/suzuki-shunsuke/gha-set-timeout-minutes@latest
```

## How to use

Please run `gha-set-timeout-minutes set` at the repository root directory.

```sh
gha-set-timeout-minutes set
```

then `gha-set-timeout-minutes` checks GitHub Actions workflows `^\.github/workflows/.*\.ya?ml$` and sets `timeout-minutes: 30` to jobs which don't have `timeout-minutes`.
Jobs which have `timeout-minutes` aren't changed.
You can specify the value of `timeout-minutes` with `-m` option.

```sh
gha-set-timeout-minutes set -m 60
```

You can specify workflow files by positional arguments.

```sh
gha-set-timeout-minutes set .github/workflows/test.yaml
```

## Tips: Fix workflows by CI

Using `gha-set-timeout-minutes` in CI, you can fix workflows automatically.
When workflow files are added or changed in a pull request, you can run `gha-set-timeout-minutes` and commit and push changes to a feature branch.

## LICENSE

[MIT](LICENSE)
