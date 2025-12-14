# ghatm

[![DeepWiki](https://img.shields.io/badge/DeepWiki-suzuki--shunsuke%2Fghatm-blue.svg?logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACwAAAAyCAYAAAAnWDnqAAAAAXNSR0IArs4c6QAAA05JREFUaEPtmUtyEzEQhtWTQyQLHNak2AB7ZnyXZMEjXMGeK/AIi+QuHrMnbChYY7MIh8g01fJoopFb0uhhEqqcbWTp06/uv1saEDv4O3n3dV60RfP947Mm9/SQc0ICFQgzfc4CYZoTPAswgSJCCUJUnAAoRHOAUOcATwbmVLWdGoH//PB8mnKqScAhsD0kYP3j/Yt5LPQe2KvcXmGvRHcDnpxfL2zOYJ1mFwrryWTz0advv1Ut4CJgf5uhDuDj5eUcAUoahrdY/56ebRWeraTjMt/00Sh3UDtjgHtQNHwcRGOC98BJEAEymycmYcWwOprTgcB6VZ5JK5TAJ+fXGLBm3FDAmn6oPPjR4rKCAoJCal2eAiQp2x0vxTPB3ALO2CRkwmDy5WohzBDwSEFKRwPbknEggCPB/imwrycgxX2NzoMCHhPkDwqYMr9tRcP5qNrMZHkVnOjRMWwLCcr8ohBVb1OMjxLwGCvjTikrsBOiA6fNyCrm8V1rP93iVPpwaE+gO0SsWmPiXB+jikdf6SizrT5qKasx5j8ABbHpFTx+vFXp9EnYQmLx02h1QTTrl6eDqxLnGjporxl3NL3agEvXdT0WmEost648sQOYAeJS9Q7bfUVoMGnjo4AZdUMQku50McDcMWcBPvr0SzbTAFDfvJqwLzgxwATnCgnp4wDl6Aa+Ax283gghmj+vj7feE2KBBRMW3FzOpLOADl0Isb5587h/U4gGvkt5v60Z1VLG8BhYjbzRwyQZemwAd6cCR5/XFWLYZRIMpX39AR0tjaGGiGzLVyhse5C9RKC6ai42ppWPKiBagOvaYk8lO7DajerabOZP46Lby5wKjw1HCRx7p9sVMOWGzb/vA1hwiWc6jm3MvQDTogQkiqIhJV0nBQBTU+3okKCFDy9WwferkHjtxib7t3xIUQtHxnIwtx4mpg26/HfwVNVDb4oI9RHmx5WGelRVlrtiw43zboCLaxv46AZeB3IlTkwouebTr1y2NjSpHz68WNFjHvupy3q8TFn3Hos2IAk4Ju5dCo8B3wP7VPr/FGaKiG+T+v+TQqIrOqMTL1VdWV1DdmcbO8KXBz6esmYWYKPwDL5b5FA1a0hwapHiom0r/cKaoqr+27/XcrS5UwSMbQAAAABJRU5ErkJggg==)](https://deepwiki.com/suzuki-shunsuke/ghatm)

`ghatm` is a command line tool setting [timeout-minutes](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#jobsjob_idtimeout-minutes) to all GitHub Actions jobs.
It finds GitHub Actions workflows and adds `timeout-minutes` to jobs which don't have the setting.
It edits workflow files while keeping YAML comments, indents, empty lines, and so on.

```console
$ ghatm set
```

```diff
diff --git a/.github/workflows/test.yaml b/.github/workflows/test.yaml
index e8c6ae7..aba3b2d 100644
--- a/.github/workflows/test.yaml
+++ b/.github/workflows/test.yaml
@@ -6,6 +6,7 @@ on: pull_request
 jobs:
   path-filter:
     # Get changed files to filter jobs
+    timeout-minutes: 30
     outputs:
       update-aqua-checksums: ${{steps.changes.outputs.update-aqua-checksums}}
       renovate-config-validator: ${{steps.changes.outputs.renovate-config-validator}}
@@ -71,6 +72,7 @@ jobs:
       contents: read
 
   build:
+    timeout-minutes: 30
     runs-on: ubuntu-latest
     permissions: {}
     steps:
```

## Motivation

- https://exercism.org/docs/building/github/gha-best-practices#h-set-timeouts-for-workflows
- [job_timeout_minutes_is_required | suzuki-shunsuke/ghalint](https://github.com/suzuki-shunsuke/ghalint/blob/main/docs/policies/012.md)
- [job_timeout_minutes_is_required | lintnet-modules/ghalint](https://github.com/lintnet-modules/ghalint/tree/main/workflow/job_timeout_minutes_is_required)

`timeout-minutes` should be set properly, but it's so bothersome to fix a lot of workflow files by hand.
`ghatm` fixes them automatically.

## How to install

`ghatm` is a single binary written in Go.
So you only need to put the executable binary into `$PATH`.

1. [Homebrew](https://brew.sh/)

```sh
brew install suzuki-shunsuke/ghatm/ghatm
```

2. [Scoop](https://scoop.sh/)

```sh
scoop bucket add suzuki-shunsuke https://github.com/suzuki-shunsuke/scoop-bucket
scoop install ghatm
```

3. [aqua](https://aquaproj.github.io/)

```sh
aqua g -i suzuki-shunsuke/ghatm
```

4. Download a prebuilt binary from [GitHub Releases](https://github.com/suzuki-shunsuke/ghatm/releases) and install it into `$PATH`

<details>
<summary>Verify downloaded assets from GitHub Releases</summary>

You can verify downloaded assets using some tools.

1. [GitHub CLI](https://cli.github.com/)
1. [slsa-verifier](https://github.com/slsa-framework/slsa-verifier)
1. [Cosign](https://github.com/sigstore/cosign)

--

1. GitHub CLI

ghatm >= v0.3.3

You can install GitHub CLI by aqua.

```sh
aqua g -i cli/cli
```

```sh
gh release download -R suzuki-shunsuke/ghatm v0.3.3 -p ghatm_darwin_arm64.tar.gz
gh attestation verify ghatm_darwin_arm64.tar.gz \
  -R suzuki-shunsuke/ghatm \
  --signer-workflow suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml
```

Output:

```
Loaded digest sha256:84298e8436f0b2c7f51cd4606848635471a11aaa03d7d0c410727630defe6b7e for file://ghatm_darwin_arm64.tar.gz
Loaded 1 attestation from GitHub API
✓ Verification succeeded!

sha256:84298e8436f0b2c7f51cd4606848635471a11aaa03d7d0c410727630defe6b7e was attested by:
REPO                                 PREDICATE_TYPE                  WORKFLOW
suzuki-shunsuke/go-release-workflow  https://slsa.dev/provenance/v1  .github/workflows/release.yaml@7f97a226912ee2978126019b1e95311d7d15c97a
```

2. slsa-verifier

You can install slsa-verifier by aqua.

```sh
aqua g -i slsa-framework/slsa-verifier
```

```sh
gh release download -R suzuki-shunsuke/ghatm v0.3.3 -p ghatm_darwin_arm64.tar.gz  -p multiple.intoto.jsonl
slsa-verifier verify-artifact ghatm_darwin_arm64.tar.gz \
  --provenance-path multiple.intoto.jsonl \
  --source-uri github.com/suzuki-shunsuke/ghatm \
  --source-tag v0.3.3
```

Output:

```
Verified signature against tlog entry index 137035428 at URL: https://rekor.sigstore.dev/api/v1/log/entries/108e9186e8c5677a421587935f03afc5f73475e880b6f05962c5be8726ccb5011b7bf62a5d2a58bb
Verified build using builder "https://github.com/slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@refs/tags/v2.0.0" at commit 1af80d4aa0b6cc45bda5677fd45202ee2b90e1fc
Verifying artifact ghatm_darwin_arm64.tar.gz: PASSED
```

3. Cosign

You can install Cosign by aqua.

```sh
aqua g -i sigstore/cosign
```

```sh
gh release download -R suzuki-shunsuke/ghatm v0.3.3
cosign verify-blob \
  --signature ghatm_0.3.3_checksums.txt.sig \
  --certificate ghatm_0.3.3_checksums.txt.pem \
  --certificate-identity-regexp 'https://github\.com/suzuki-shunsuke/go-release-workflow/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  ghatm_0.3.3_checksums.txt
```

Output:

```
Verified OK
```

After verifying the checksum, verify the artifact.

```sh
cat ghatm_0.3.3_checksums.txt | sha256sum -c --ignore-missing
```

</details>

5. Go

```sh
go install github.com/suzuki-shunsuke/ghatm/cmd/ghatm@latest
```

## How to use

Please run `ghatm set` on the repository root directory.

```sh
ghatm set
```

Then `ghatm` checks GitHub Actions workflows `^\.github/workflows/.*\.ya?ml$` and sets `timeout-minutes: 30` to jobs not having `timeout-minutes`.
Jobs with `timeout-minutes` aren't changed.
You can specify the value of `timeout-minutes` with `-t` option.

```sh
ghatm set -t 60
```

You can specify workflow files by positional arguments.

```sh
ghatm set .github/workflows/test.yaml
```

### Decide `timeout-minutes` based on each job's past execution times

```sh
ghatm set -auto [-repo <repository>] [-size <the number of sample data>]
```

ghatm >= v0.3.2 [#68](https://github.com/suzuki-shunsuke/ghatm/issues/68) [#70](https://github.com/suzuki-shunsuke/ghatm/pull/70)

> [!warning]
> The feature doesn't support workflows using `workflow_call`.

If the `-auto` option is used, ghatm calls GitHub API to get each job's past execution times and decide appropriate `timeout-minutes`.
This feature requires a GitHub access token with the `actions:read` permission.
You have to set the access token to the environment variable `GITHUB_TOKEN` or `GHATM_GITHUB_TOKEN`.

GitHub API:

- [List workflow runs for a workflow](https://docs.github.com/en/rest/actions/workflow-runs?apiVersion=2022-11-28#list-workflow-runs-for-a-workflow)
- [List jobs for a workflow run](https://docs.github.com/en/rest/actions/workflow-jobs#list-jobs-for-a-workflow-run)

ghatm takes 30 jobs by job to decide `timeout-minutes`.
You can change the number of jobs by the `-size` option.

```
max(job execution times) + 10
```

## Tips: Fix workflows by CI

Using `ghatm` in CI, you can fix workflows automatically.
When workflow files are added or changed in a pull request, you can run `ghatm` and commit and push changes to a feature branch.

## LICENSE

[MIT](LICENSE)
