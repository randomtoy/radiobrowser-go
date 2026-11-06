# Changelog

## [0.3.0] - 2026-11-06

This release completes the current, non-deprecated Radio Browser JSON API coverage
and fixes compatibility with live API responses. It includes a source-breaking
change to station coordinates; follow the migration instructions below when
upgrading from v0.1.x or v0.2.x.

### Breaking changes

- `Station.GeoLat` and `Station.GeoLong` are now `*float64` instead of `float64`.
  Unknown coordinates remain `nil`; known zero coordinates remain pointers to zero.
  Code that reads, compares or initializes these fields must handle pointers.

### Added

- `Stats`, `Servers` and `Config` for `/json/stats`, `/json/servers` and
  `/json/config`, with `ServerStats`, `ServerMirror` and `ServerConfig` models.
- `GeoLat`, `GeoLong` and `GeoDistance` advanced-search options. Pointers preserve
  explicit zero values; distance and radius are measured in meters.
- `Limit` for countries, codecs, states, languages and tags queries.
- `LimitSet` across paginated options to distinguish an explicit `limit=0` from
  the existing server-default behavior.
- `ISO3166_2` and `LanguageCodes` in `AddStationRequest`.
- `Station.ISO3166_2`, `LastLocalCheckTimeISO`, `SSLError` and nullable `GeoDistance`.
- `Message`, `StationUUID`, `Name` and `URL` in `StationClickCounterResponse`.
- Geographic validation: latitude must be within [-90, 90], longitude within
  [-180, 180], and radius must be non-negative. All supplied values must be finite.
  Coordinate validation also applies to station creation.
- A server information example, geographic search examples and package documentation.

### Fixed

- `StationsByClicks` now calls `/json/stations/topclick`.
- `StationsByURL` now calls `/json/stations/byurl` with the stream URL in the
  `url` query parameter, preserving the URL's own query string.
- `StationClickCounter` now calls `/json/url/{stationuuid}` and retains the full
  response. Its `OK` flag accepts both JSON booleans and documented quoted booleans.
- Path segments are encoded exactly once, including spaces, Unicode, literal
  percent signs, escaped slashes and other reserved characters.
- Country-only `States` queries preserve the empty region filter segment.
- All five metadata models accept numeric and legacy quoted `stationcount` values
  while retaining ordinary public `int` fields.
- README health-check examples use `ChecksOptions`; server descriptions distinguish
  API mirrors from audio streaming servers.
- Popularity examples label clicks as a 24-hour metric, metadata examples request
  bounded result sets, and truncation preserves Unicode characters.
- Example newline formatting no longer fails `go vet` or `go test ./...`.

### Deprecated

These fields remain available for compatibility:

| Field | Replacement |
|-------|-------------|
| `Station.Country` | `Station.CountryCode` |
| `AddStationRequest.State` | `AddStationRequest.ISO3166_2` |
| `AddStationRequest.Language` | `AddStationRequest.LanguageCodes` |

Country-name search endpoints remain supported. The deprecated `/countrycodes`
endpoint, non-JSON exports and Prometheus `/metrics` are outside this release's
JSON API scope.

### Migration from v0.1.x and v0.2.x

1. Update the dependency and refresh your module files:

   ```bash
   go get github.com/randomtoy/radiobrowser-go@v0.3.0
   go mod tidy
   ```

2. Check station coordinates for `nil` before using their values. When constructing
   a station, pass pointers to coordinate variables:

   ```go
   lat, lon := 52.37, 4.90
   station := rb.Station{GeoLat: &lat, GeoLong: &lon}

   if station.GeoLat != nil && station.GeoLong != nil {
       fmt.Printf("%.4f, %.4f\n", *station.GeoLat, *station.GeoLong)
   }
   ```

   Do not replace unknown coordinates with zero: `(0, 0)` is a valid location.
   `Station.GeoDistance` is also nullable.

3. Prefer code-based fields when creating stations. `LanguageCodes` contains
   comma-separated ISO 639-2/B codes; `ISO3166_2` identifies a country subdivision:

   ```go
   request := rb.AddStationRequest{
       Name:          "My Radio",
       URL:           "https://example.com/live",
       CountryCode:   "NL",
       ISO3166_2:     "NL-NH",
       LanguageCodes: "dut,eng",
   }
   ```

4. Existing pagination calls keep their behavior: `Limit: 0` omits the parameter.
   Only set `LimitSet: true` when you intend to send an explicit zero limit:

   ```go
   options := rb.StationsOptions{Limit: 0, LimitSet: true}
   ```

   Valid zero coordinates are sent through pointers. Invalid geographic values now
   return an error before an HTTP request, so check errors from `StationSearch`
   and `AddStation`.

5. `DefaultUserAgent()` and `WithUserAgent` are unchanged from v0.2.0. Consumers
   built with the v0.3.0 dependency automatically use `radiobrowser-go/0.3.0`.
   Unversioned builds and local module replacements still use `radiobrowser-go/dev`.
   No source-level version override is needed.

### Validation

- `go test ./...`
- `go test -race ./...`
- `go vet ./...`
- `go build ./...`
- All 17 README Go snippets compile; all five read-only examples were exercised
  against the public API. Station creation, voting and play registration were
  verified against local HTTP test servers.

### Related work

- [#3](https://github.com/randomtoy/radiobrowser-go/issues/3): routes and URL encoding
  ([PR #8](https://github.com/randomtoy/radiobrowser-go/pull/8)).
- [#4](https://github.com/randomtoy/radiobrowser-go/issues/4): metadata decoding and
  response models ([PR #9](https://github.com/randomtoy/radiobrowser-go/pull/9)).
- [#5](https://github.com/randomtoy/radiobrowser-go/issues/5): search, creation and
  pagination options ([PR #10](https://github.com/randomtoy/radiobrowser-go/pull/10)).
- [#6](https://github.com/randomtoy/radiobrowser-go/issues/6): server JSON endpoints
  ([PR #11](https://github.com/randomtoy/radiobrowser-go/pull/11)).
- [#7](https://github.com/randomtoy/radiobrowser-go/issues/7): documentation and examples
  ([PR #13](https://github.com/randomtoy/radiobrowser-go/pull/13)).

[0.3.0]: https://github.com/randomtoy/radiobrowser-go/compare/v0.2.0...v0.3.0
