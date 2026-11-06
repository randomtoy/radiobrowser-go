# radiobrowser-go

[![Go Reference](https://pkg.go.dev/badge/github.com/randomtoy/radiobrowser-go.svg)](https://pkg.go.dev/github.com/randomtoy/radiobrowser-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A zero-dependency Go client for the [Radio Browser JSON API](https://de1.api.radio-browser.info/), a community database of internet radio stations. Requires Go 1.22 or later.

## Features

- Station discovery by name, country code, language, codec, tags, UUID and stream URL.
- Advanced filters, geographic radius searches, pagination and popularity rankings.
- Countries, codecs, regions, languages and tags with station counts.
- Availability checks, check steps, click history and station change history.
- Voting, play registration and station creation with current subdivision/language codes.
- API server statistics, mirror lists and configuration.
- Context cancellation and application-specific User-Agent headers.

The library covers the current, non-deprecated JSON endpoints documented by Radio Browser 0.7.45. Non-JSON exports, Prometheus `/metrics` and the deprecated `/countrycodes` endpoint are outside its scope.

## Installation

```bash
go get github.com/randomtoy/radiobrowser-go
```

## Quick start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    rb "github.com/randomtoy/radiobrowser-go"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()
    ctx = rb.WithUserAgent(ctx, "MyRadio/1.0")

    stations, err := rb.StationsSearch(ctx, "https://de1.api.radio-browser.info",
        rb.StationByCountryCodeExact, "DE", rb.StationsOptions{
            Limit: 10, Order: rb.StationOrderVotes, Reverse: true, HideBroken: true,
        })
    if err != nil {
        log.Fatal(err)
    }
    for _, station := range stations {
        fmt.Printf("%s (%s): %d votes\n", station.Name, station.CountryCode, station.Votes)
    }
}
```

Every API function takes a context and a base URL. Set a deadline on the context to bound request duration; the shared HTTP client has no default timeout. HTTP, network and JSON decoding failures are returned as errors. For action responses, also inspect `OK` and `Message`: an HTTP-successful response can report an unsuccessful operation.

### User-Agent

Requests use `rb.DefaultUserAgent` (`radiobrowser-go/0.2.0`) unless overridden. Radio Browser asks clients to identify their application and version:

```go
ctx := rb.WithUserAgent(context.Background(), "MyRadio/1.2")
```

The setting is inherited by derived contexts. Whitespace is trimmed; an empty value restores the default.

## Search and pagination

### Advanced filters

```go
stations, err := rb.StationSearch(ctx, baseURL, rb.StationSearchOptions{
    CountryCode: "US",
    TagList:     []string{"jazz", "instrumental"}, // Every tag must match.
    Codec:       "MP3",
    BitrateMin:  128,
    Limit:       20,
    Order:       rb.StationOrderClickCount,
    Reverse:     true,
    HideBroken:  true,
})
```

`HasGeoInfo`, `HasExtendedInfo` and `IsHTTPS` are `*bool`: `nil` leaves the filter unset, while pointers to `true` or `false` select either side. `NameExact`, `CountryExact`, `StateExact`, `LanguageExact` and `TagExact` enable exact matches for their corresponding string filters.

### Geographic search

```go
lat, lon, radius := 52.37, 4.90, 50000.0
stations, err := rb.StationSearch(ctx, baseURL, rb.StationSearchOptions{
    GeoLat: &lat, GeoLong: &lon, GeoDistance: &radius,
    Limit: 10, HideBroken: true,
})
```

Coordinates and radius use `*float64`, so zero remains a valid explicit value. Latitude must be between -90 and 90, longitude between -180 and 180, and radius must be non-negative and finite. Radius and returned `Station.GeoDistance` are measured in meters. Provide both coordinates for a meaningful distance search.

### Limits and offsets

All paginated option types support `Limit` and `LimitSet`. A positive `Limit` is sent automatically. The zero-value options omit the limit and use the server default; to explicitly request zero rows, set `LimitSet: true`:

```go
stations, err := rb.Stations(ctx, baseURL, rb.StationsOptions{
    Limit: 0, LimitSet: true,
})
```

For ordinary paging use, for example, `Limit: 20, Offset: 40`. Metadata options also support `Order`, `Reverse`, `HideBroken`, `Offset` and `Limit`. Change, check and click history use their UUID cursor fields rather than offsets. Ranking endpoints apply their own order; use `Stations` or `StationSearch` for custom `Order`/`Reverse` sorting.

### UUID and URL lookup

```go
stations, err := rb.StationsByUUID(ctx, baseURL, []string{stationUUID})
```

```go
stations, err := rb.StationsByURL(ctx, baseURL,
    "https://example.com/live?channel=jazz", rb.StationsOptions{})
```

URLs are sent as query parameters, preserving their own query strings. Text search and metadata filters support spaces, Unicode and reserved characters.

## Metadata and station health

```go
countries, err := rb.Countries(ctx, baseURL, "", rb.CountriesOptions{
    Order: "stationcount", Reverse: true, Limit: 10,
})
```

```go
states, err := rb.States(ctx, baseURL, "Germany", "", rb.StatesOptions{Limit: 10})
```

For `States`, `country` is a country name; the independent `filter` argument matches region names. An empty filter lists all regions of the selected country. Countries expose their ISO country code through `ISO3166_1`. All five metadata models decode both current numeric station counts and legacy quoted counts into `int`.

```go
checks, err := rb.StationCheckResults(ctx, baseURL, stationUUID, rb.ChecksOptions{Limit: 5})
```

```go
steps, err := rb.StationCheckSteps(ctx, baseURL, []string{stationUUID})
```

An empty station UUID retrieves the latest checks across stations. Supplying a UUID selects that station's check history. `StationClicks` and `StationsChanged` similarly accept an empty UUID for a global feed.

## Station actions

Register a play when the user starts playback. The API counts at most one play per IP address and station per day; the response includes the station name and playable URL:

```go
play, err := rb.StationClickCounter(ctx, baseURL, stationUUID)
if err != nil {
    log.Fatal(err)
}
if !play.OK {
    log.Fatal(play.Message)
}
fmt.Println(play.URL)
```

Vote with `VoteForStation`. Create a station with the current code-based fields:

```go
added, err := rb.AddStation(ctx, baseURL, rb.AddStationRequest{
    Name:          "My Radio",
    URL:           "https://example.com/live",
    Homepage:      "https://example.com",
    CountryCode:   "NL",
    ISO3166_2:     "NL-NH",
    LanguageCodes: "dut,eng",
    Tags:          "jazz,instrumental",
})
```

`LanguageCodes` uses comma-separated ISO 639-2/B codes. `ISO3166_2` identifies a country subdivision. The legacy creation fields `State` and `Language` remain available but are deprecated. `GeoLat` and `GeoLong` are optional pointers, with the same coordinate validation as search.

## API servers

```go
stats, err := rb.Stats(ctx, baseURL)
```

```go
config, err := rb.Config(ctx, baseURL)
```

```go
mirrors, err := rb.Servers(ctx, baseURL)
```

`Servers` lists Radio Browser API mirrors with their hostname and IPv4/IPv6 address. Radio Browser recommends DNS lookup of `all.api.radio-browser.info`, followed by reverse lookup, when client-side DNS is available; the HTTP endpoint is its fallback. The library accepts the selected server URL explicitly and does not automatically discover or retry mirrors.

`StreamingServers` is a different endpoint: it lists audio streaming servers such as Icecast, including their administrative/status metadata. It does not list API mirrors.

## API reference

All paths below use the `/json` prefix. Optional filters, station UUIDs and pagination parameters are omitted for brevity.

| Function | Endpoint / purpose |
|----------|--------------------|
| `Stations` | `/stations` — all stations |
| `StationsSearch` | `/stations/{kind}/{term}` — all 14 documented simple search variants |
| `StationSearch` | `/stations/search` — advanced filters and geographic search |
| `StationsByUUID` | `/stations/byuuid?uuids=...` |
| `StationsByURL` | `/stations/byurl?url=...` |
| `StationsByClicks` | `/stations/topclick` — most clicks in the last 24 hours |
| `StationsByVotes` | `/stations/topvote` |
| `StationsByRecentClick` | `/stations/lastclick` |
| `StationsByRecentlyChanged` | `/stations/lastchange` |
| `StationsChanged` | `/stations/changed` — old station versions |
| `StationsBroken` | `/stations/broken` |
| `Countries` | `/countries` |
| `Codecs` | `/codecs` |
| `States` | `/states` |
| `Languages` | `/languages` |
| `Tags` | `/tags` |
| `StreamingServers` | `/streamingservers` — audio streaming servers |
| `StationCheckResults` | `/checks` |
| `StationCheckSteps` | `/checksteps?uuids=...` |
| `StationClicks` | `/clicks` — play history |
| `StationClickCounter` | `/url/{uuid}` — register playback and retrieve the stream URL |
| `VoteForStation` | `/vote/{uuid}` |
| `AddStation` | `/add` — HTTP POST form, JSON response |
| `Stats` | `/stats` — `ServerStats` |
| `Servers` | `/servers` — `[]ServerMirror` |
| `Config` | `/config` — `ServerConfig` |

Search order constants cover all documented sort keys, including random order. `ClickCount` covers the last 24 hours; `ClickTrend` compares the last two days.

## Migration from v0.1.x

`Station.GeoLat` and `Station.GeoLong` are now `*float64`. Check for `nil` before dereferencing; unknown coordinates are no longer confused with `(0, 0)`:

```go
if station.GeoLat != nil && station.GeoLong != nil {
    fmt.Printf("%.4f, %.4f\n", *station.GeoLat, *station.GeoLong)
}
```

Station responses also expose `ISO3166_2`, `LastLocalCheckTimeISO`, `SSLError` and nullable `GeoDistance`. `Country` remains available but is deprecated in favor of `CountryCode`; country-name search endpoints remain supported. Existing `Limit: 0` behavior is preserved unless `LimitSet` is enabled. Metadata `StationCount` fields remain ordinary `int` values.

## Examples and validation

| Example | Description |
|---------|-------------|
| [basic](./examples/basic) | Country-code search |
| [search](./examples/search) | Advanced filters, text search and geographic radius |
| [metadata](./examples/metadata) | Bounded countries, tags, codecs, languages and regions |
| [topstations](./examples/topstations) | Votes, 24-hour clicks, recent activity and trends |
| [server](./examples/server) | Server statistics and configuration; optional HTTP mirror lookup |

Run from the repository root:

```bash
go run ./examples/basic
go run ./examples/search
go run ./examples/metadata
go run ./examples/topstations
go run ./examples/server
```

Examples only read API data. Tests exercise station creation, voting and play registration against local HTTP servers. The regular test suite does not depend on the public API:

```bash
go test ./...
go test -race ./...
go vet ./...
go build ./...
```

## License

[MIT](./LICENSE).
