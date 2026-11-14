# quakes

Prints the past day's earthquakes from the [USGS earthquake feed](https://earthquake.usgs.gov/earthquakes/feed/v1.0/geojson.php).

## Install

```sh
go install github.com/mikemadden42/quakes@latest
```

Prebuilt binaries for Linux, macOS, and Windows are attached to [releases](https://github.com/mikemadden42/quakes/releases).

## Usage

```sh
quakes
```

Each line shows the time (local zone), event type, magnitude, location, and depth. An unknown magnitude or depth prints as `?`.

```
USGS All Earthquakes, Past Day
2026-11-14 09:42:36 CST  earthquake   1.5  101 km E of Cordova, Alaska  28.1 km
2026-11-14 09:39:37 CST  earthquake   1.3  45 km WNW of Valdez, Alaska  5.0 km
2026-11-14 09:37:02 CST  earthquake   1.3  2 km NW of The Geysers, CA  1.5 km
```

## Development

| Command | What it does |
|---|---|
| `go test ./...` | Run tests |
| `make check` | Lint with golangci-lint |
| `make build` | Build binaries for all platforms via goreleaser (output in `dist/`) |
| `make snapshot` | Build binaries plus release archives and checksums locally, without publishing |

## License

MIT, see [LICENSE](LICENSE).
