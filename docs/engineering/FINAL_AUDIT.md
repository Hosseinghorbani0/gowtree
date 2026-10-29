# gowtree Engineering Upgrade — Final Evidence Report

## Scope

This continuation preserved the existing staged and unstaged worktree changes and addressed the highest-priority unfinished correctness, error-handling, traversal, performance, and release-versioning gaps. No release was published and no tag was pushed.

## Verified fixes

| Area | Evidence |
|---|---|
| Explicit config path | `internal/app.Run` rejects a missing `--config` file; covered by `TestRunRejectsExplicitlyMissingConfig`. |
| YAML errors | `internal/config.Load` returns malformed YAML errors; covered by `internal/config/config_test.go`. |
| CLI validation | Color, charset, sort, time, depth, mode conflicts, output conflicts, and invalid ignore globs are rejected by `validateOptions`. |
| Symlinks | `internal/tree.classifyEntry` distinguishes symlinks, broken symlinks, and cycles. Existing regression tests cover JSON behavior with and without `-l`. |
| Traversal errors | `Print`, `BuildJSONStrict`, and `CountItemsStrict` return filesystem errors to `internal/app`. |
| Output errors | Text statistics/root output and Markdown/HTML renderers now check writer errors. Renderer failures are covered in `internal/output/render_test.go`. |
| CLI contracts | `internal/app/cli_contract_test.go` verifies stdout/stderr separation, invalid-argument and invalid-path exit codes, JSON decoding, and closed-output failure behavior. The contract tests pass repeatedly with `-count=10`. |
| Partial writes | `internal/output/render_test.go` verifies renderer propagation of `io.ErrShortWrite`. |
| Metadata defaults | `internal/metadata/metadata_test.go` verifies non-empty metadata and mutable linker-injection behavior. |
| HTML safety | HTML filenames are escaped; regression coverage verifies `<` and `&` escaping. |
| Count semantics | `CountItemsStrict` excludes the root and applies the same filtering/depth/symlink rules used by `Print`. |
| `--du` repeated scans | `Options.SizeCache` caches directory sizes during one run and is shared by text/JSON traversal. |
| Release version injection | `metadata.Version` is now a variable, allowing `-ldflags -X`. A real binary built with `v9.8.7` printed `gowtree v9.8.7`. |
| Installer naming | Release workflow passes the tag-derived version to `build-installer.ps1`; Inno Setup receives `/DMyAppVersion=...`, and artifact names use the same value. |

## Tests and static analysis

Executed in `D:\golang`:

- `go test ./...` — passed.
- `go vet ./...` — passed.
- `go test -cover ./...` — passed:
  - `internal/format`: 100.0%
  - `internal/config`: 61.5%
  - `internal/tree`: 42.7%
  - `internal/output`: 22.1%
  - `internal/platform`: 83.3% (Windows-specific tests executed on this Windows runner)
  - `internal/style`: 68.6%
  - `internal/app`: 46.8% after CLI contract tests
  - `internal/metadata`: no statements reported; defaults and linker-compatible mutability are tested.
  - `cmd/gowtree`: 0.0% reported; executable behavior is covered by build/smoke tests.
- `git diff --check` — passed.
- `go test -race ./...` — not runnable on this Windows environment because `CGO_ENABLED=0`; Go reported that race testing requires CGO.
- `actionlint` — unavailable in the environment; YAML was parsed successfully with Python/PyYAML, which is not a replacement for actionlint.
- `govulncheck` — executed with `go run golang.org/x/vuln/cmd/govulncheck@v1.1.4 ./...`; no vulnerabilities were reported in reachable code. The tool also reported 3 vulnerabilities in imported packages and 21 in required modules that it did not identify as reachable from this code; these still require dependency-owner review.
- `go test ./internal/app -run 'TestRunCLIContracts|TestRunReportsOutputWriterFailure' -count=10` — passed, including independent stream capture and JSON decoding.
- `go test ./internal/output -run 'TestRender' -count=10` — passed.
- `go build ./cmd/gowtree` and `--version` smoke test — passed (`gowtree v1.4.0`).
- `go build -trimpath ./cmd/gowtree` and `--version` smoke test — passed.
- Linker injection smoke test — passed with `-X ...metadata.Version=v7.6.5`; binary printed `gowtree v7.6.5`.
- Workflow and Dependabot YAML parsing with Python/PyYAML — passed.

## Benchmark evidence

Command:

```text
go test ./internal/tree -run '^$' -bench BenchmarkDirectorySizeRepeatedScans -benchmem -count=1
```

Fixture: 8 directories, each containing 32 files of 256 bytes, on Windows/amd64, Go 1.26.4.

| Case | ns/op | B/op | allocs/op |
|---|---:|---:|---:|
| uncached repeated scan | 5,193,550 | 124,303 | 1,149 |
| cached repeated lookup | 1,224 | 480 | 3 |

This is an actual benchmark for the local fixture, not a general production speed-up claim. The cached case intentionally measures repeated lookups after the first calculation; it demonstrates the avoided filesystem work, not an end-to-end CLI benchmark.

## Release and security validation

- The release workflow is tag-triggered, validates `vMAJOR.MINOR.PATCH`, and passes the normalized version to the installer script.
- Binary metadata injection was executed locally and verified with `--version`.
- Inno Setup was not executed locally because the required installer compiler was not available.
- The release workflow now verifies the downloaded Inno Setup `6.7.3` executable against the SHA-256 digest published in the official `jrsoftware/issrc` GitHub release metadata: `9c73c3bae7ed48d44112a0f48e66742c00090bdb5bef71d9d3c056c66e97b732`. The `.issig` signature is not independently verified because a trusted public key is not maintained by this repository.
- The workflow also rejects malformed configured digests, missing downloaded files, hash mismatches, and download failures before executing the installer.
- The official GitHub release API was queried for the exact `innosetup-6.7.3.exe` asset and reported the same SHA-256 digest used by the workflow; the workflow itself was not run on GitHub.
- GitHub Actions, CodeQL, artifact publication, and rollback behavior were not executed here. `govulncheck` was executed locally as recorded above; its GitHub workflow execution remains unverified.
- `gh workflow list` showed only the historical remote `CI/CD` workflow. The latest inspected successful remote run was `35518821149` at an older commit SHA, using the old unpinned workflow; it is not evidence that the current local workflows passed. The local workflows cannot be run remotely without pushing them, which was not authorized.

## Files changed in this continuation

- `internal/app/app.go`
- `internal/app/app_test.go`
- `internal/app/validation.go`
- `internal/app/validation_test.go`
- `internal/config/config_test.go`
- `internal/metadata/metadata.go`
- `internal/output/render.go`
- `internal/output/render_test.go`
- `internal/tree/tree.go`
- `internal/tree/tree_test.go`
- `scripts/build-installer.ps1`
- `installer/gowtree.iss`
- `.github/workflows/release.yml`
- `README.md`
- `docs/engineering/PLAN.md`
- `docs/engineering/FINAL_AUDIT.md`
- `internal/style/style_test.go`
- `internal/platform/filetime_windows_test.go`
- `internal/platform/hidden_windows_test.go`
- `internal/app/cli_contract_test.go`
- `internal/metadata/metadata_test.go`

The worktree also contains earlier user/agent changes in CI, assets, installer files, and documentation; they were preserved.

## Remaining bugs and blockers

1. **Medium — release-tool signature key management.** `release.yml` verifies the official GitHub release SHA-256 digest for Inno Setup `6.7.3`, but does not independently verify the accompanying `.issig` signature because a trusted public key is not maintained by this repository.
2. **Low — legacy `DirSize` API suppresses walk errors.** Production traversal now uses `DirSizeWithError`/the cached strict path, but the compatibility `DirSize` wrapper still returns only `int64` for callers that use the old API.
3. **Low — compatibility sorting wrapper.** Production `Print` uses `SortEntriesStrict`; the legacy `SortEntries` wrapper retains a fallback-only signature for compatibility.
4. **Medium — low test coverage outside core traversal.** `internal/output` branches remain at 22.1%; `cmd/gowtree` has no package-level tests, although executable smoke tests and `internal/app` contract tests cover entrypoint behavior.
5. **Medium — race validation unavailable locally.** A CGO-enabled Windows or Linux CI runner is required for actual race evidence.
6. **Low — renderer model.** Markdown and HTML intentionally render only `directory` and `file` nodes; symlink-specific JSON node types are not visibly labeled in those formats.

## Prioritized next steps

1. Add and test trusted Inno Setup `.issig` verification in the release workflow.
2. Deprecate the compatibility `DirSize(path, verbose) int64` wrapper in favor of `DirSizeWithError(path) (int64, error)` for future callers.
3. Define a diagnostics policy for `DirEntry.Info()` failures and test inaccessible/restricted directories on Windows.
4. Push the reviewed local workflows through the normal review process, then execute `actionlint`, `govulncheck`, race tests, installer compilation, and release dry runs on CI; this was not performed from the current worktree.
5. Expand app/output/platform coverage and add end-to-end contract tests for JSON, Markdown, HTML, and text output.
6. Benchmark complete `--du` CLI runs on flat, deep, filtered, and symlink-heavy fixtures before making further optimization claims.

Passing unit tests do not establish production readiness until the release supply chain, race testing, installer build, and filesystem permission cases are executed in CI.
