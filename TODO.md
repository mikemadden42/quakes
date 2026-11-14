# TODO

## quakes.go

### Bugs / robustness
- [x] Add an HTTP timeout — `http.Get` uses the default client, which never times out (`&http.Client{Timeout: 15 * time.Second}`)
- [x] Close the response body (`defer response.Body.Close()`)
- [x] Check `response.StatusCode != http.StatusOK` before decoding, so a 5xx/HTML error page doesn't end in a JSON panic
- [x] Guard `quake.Geometry.Coordinates[2]` with a `len()` check to avoid an index-out-of-range panic
- [x] Use `*float64` for `Mag` so a `null` magnitude isn't silently printed as `0`

### Cleanup
- [x] Pick one error-handling style (`log.Fatal`) and drop `checkErr` / `fmt.Print` + `os.Exit` / `panic`
- [x] Replace `io.ReadAll` + `json.Unmarshal` with `json.NewDecoder(response.Body).Decode(...)`
- [x] Use `time.UnixMilli(quake.Properties.Time)` instead of `time.Unix(t/1000, 0)`
- [x] Format output — depth prints as e.g. `1.49000000953674`; use `fmt.Printf` with `%.1f` or `text/tabwriter`
- [x] Trim the struct to the ~5 fields actually used; `interface{}` → `any`
- [x] Fix lint findings: add a package comment, rename `Ids` → `IDs`

## Tooling / config
- [x] Migrate `.golangci.yml` to v2 schema (`golangci-lint config verify` currently fails):
  - [x] `version: "2"` (string, not number)
  - [x] `disable-all: true` → `linters.default: none`
  - [x] `linters-settings` → `linters.settings`
  - [x] `issues.exclude-dirs` → `linters.exclusions.paths`
  - [x] `goimports.local-prefixes` → `formatters.settings.goimports.local-prefixes` (list)
  - [x] Consider enabling `bodyclose` and `staticcheck`
- [x] Fix `.goreleaser.yaml` (`goreleaser check` currently fails):
  - [x] Add `version: 2` header
  - [x] `snapshot.name_template` → `snapshot.version_template`
  - [x] Remove the no-op `go generate ./...` hook
- [x] Update CI matrix to supported Go versions (1.26.x, 1.27.x)
- [x] Bump `go.mod` from `go 1.18` to the minimum version tested
- [ ] Decide on `make release` — it passes `--snapshot`, so it never publishes and nearly duplicates `make build`
- [ ] Remove `Jenkinsfile` unless Jenkins is still in use (duplicates GitHub Actions, builds outside goreleaser)
- [x] Add `dist/` and the `quakes` binary to `.gitignore`

## Project hygiene
- [x] Add tests — split fetch/print into functions taking `io.Reader`/`io.Writer` and test against a saved GeoJSON fixture
- [x] Flesh out `README.md`: what it does, install/run, sample output
- [ ] Add flags for feed selection (`all_hour`, `2.5_day`, `significant_week`, …) and/or minimum magnitude
