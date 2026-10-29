# gowtree Engineering Hardening Plan

## Baseline and current status

- `go test ./...`: passed
- `go vet ./...`: passed
- `go test -cover ./...`: passed; coverage is concentrated in `internal/format` and `internal/tree`
- `go test -race ./...`: not runnable in the current environment because CGO is disabled
- `govulncheck@v1.1.4 ./...`: executed locally; no reachable-code vulnerabilities reported, with non-reachable dependency advisories also reported
- Windows CLI smoke tests: symlink handling was exercised; JSON without `-l` reports symlinks as nodes and `-l` detects parent cycles
- Explicit missing config, malformed YAML, invalid enums, and invalid ignore globs now return non-zero status
- Markdown/HTML renderers and strict JSON traversal now propagate writer/filesystem errors
- `CountItemsStrict` follows the same filtering, depth, and symlink-pruning rules as `Print`
- `--du` uses a per-run directory-size cache; benchmark evidence is recorded in `FINAL_AUDIT.md`

## Execution order

1. P0 correctness: configuration and CLI validation, error propagation, symlink edge cases
2. P1 architecture: shared traversal semantics and deterministic output
3. P1 performance: measure and reduce repeated `--du` scans
4. P1 quality: error-path, integration, and filesystem regression tests
5. P1/P2 security and supply-chain verification
6. P1 release versioning and packaging consistency
7. P2 UX and documentation alignment
8. Final verification and evidence report

## Remaining work

- Replace the legacy error-suppressing `DirSize` API with an error-returning variant.
- Decide whether `SortEntries` metadata failures should be fatal or explicitly surfaced as diagnostics.
- Maintain the trusted Inno Setup public key in-repository before adding `.issig` verification; the current workflow verifies the official GitHub release SHA-256 digest.
- Run `actionlint`, race tests with CGO enabled, and the installer workflow in CI. `govulncheck@v1.1.4` has now run locally; review its reported non-reachable dependency advisories and repeat it in CI.
- The GitHub repository currently exposes only the historical `CI/CD` workflow; the local `ci.yml`, `security.yml`, and release changes are not remotely executable until an authorized push occurs. No push or workflow dispatch was performed.

## Constraints

- Preserve unrelated user changes.
- Preserve existing CLI flags and output formats unless a compatibility decision is documented.
- Do not publish releases or modify the read-only task attachment.
