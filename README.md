# gha-set-timeout-minutes

Set [timeout-minutes](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#jobsjob_idtimeout-minutes) to all GitHub Actions jobs

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

## LICENSE

[MIT](LICENSE)
