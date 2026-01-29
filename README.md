# radiobrowser-go

[![Go Reference](https://pkg.go.dev/badge/github.com/randomtoy/radiobrowser-go.svg)](https://pkg.go.dev/github.com/randomtoy/radiobrowser-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A simple, zero-dependency Go client library for the [Radio Browser API](https://www.radio-browser.info/) - a free and open-source database of internet radio stations.

## Features

- **Station Discovery** - Search and retrieve radio stations by name, country, language, codec, tag, and more
- **Multiple Search Strategies** - Simple search by category, advanced multi-filter search, or lookup by UUID/URL
- **Popularity Metrics** - Get top stations by votes, clicks, recent activity, or trending
- **Metadata Access** - Retrieve lists of countries, languages, codecs, states, and tags with station counts
- **Station Health Checks** - Access availability check results and detailed check steps
- **User Interactions** - Vote for stations, track clicks, add new stations to the database
- **Zero Dependencies** - Uses only Go standard library
- **Context Support** - All API functions accept `context.Context` for timeout and cancellation

## Installation

```bash
go get github.com/randomtoy/radiobrowser-go
```

Requires Go 1.22 or later.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    rb "github.com/randomtoy/radiobrowser-go"
)

func main() {
    ctx := context.Background()
    baseURL := "https://de1.api.radio-browser.info"

    // Search stations by country
    stations, err := rb.StationsSearch(ctx, baseURL, rb.StationByCountry, "Germany", rb.StationsOptions{
        Limit: 10,
        Order: rb.StationOrderVotes,
    })
    if err != nil {
        log.Fatal(err)
    }

    for _, s := range stations {
        fmt.Printf("%s - %s (%d votes)\n", s.Name, s.Country, s.Votes)
    }
}
```

## Examples

See the [examples](./examples) directory for complete working examples:

| Example | Description |
|---------|-------------|
| [basic](./examples/basic) | Simple station search by country |
| [search](./examples/search) | Advanced search with multiple filters |
| [metadata](./examples/metadata) | Fetching countries, languages, tags, codecs |
| [topstations](./examples/topstations) | Getting popular and trending stations |

Run an example:

```bash
cd examples/basic
go run main.go
```

## Usage

### Advanced Search with Multiple Filters

```go
stations, err := rb.StationSearch(ctx, baseURL, rb.StationSearchOptions{
    Country:    "United States",
    Tag:        "jazz",
    BitrateMin: 128,
    Codec:      "MP3",
    Limit:      20,
    Order:      rb.StationOrderClickCount,
    HideBroken: true,
})
```

### Get Top Stations

```go
// Top voted stations
topVoted, err := rb.StationsByVotes(ctx, baseURL, rb.StationsOptions{Limit: 10})

// Most clicked stations
topClicked, err := rb.StationsByClicks(ctx, baseURL, rb.StationsOptions{Limit: 10})

// Recently played stations
recentlyPlayed, err := rb.StationsByRecentClick(ctx, baseURL, rb.StationsOptions{Limit: 10})
```

### Lookup Stations by UUID

```go
stations, err := rb.StationsByUUID(ctx, baseURL, []string{
    "96202c3a-0601-11e8-ae97-52543be04c81",
    "960291ca-0601-11e8-ae97-52543be04c81",
})
```

### Get Metadata

```go
// Get all countries with station counts
countries, err := rb.Countries(ctx, baseURL, "", rb.CountriesOptions{})

// Get available codecs
codecs, err := rb.Codecs(ctx, baseURL, "", rb.CodecsOptions{})

// Get languages
languages, err := rb.Languages(ctx, baseURL, "", rb.LanguagesOptions{})

// Get popular tags
tags, err := rb.Tags(ctx, baseURL, "", rb.TagsOptions{
    Order:   "stationcount",
    Reverse: true,
})
```

### Vote for a Station

```go
resp, err := rb.VoteForStation(ctx, baseURL, "96202c3a-0601-11e8-ae97-52543be04c81")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Vote result: %v - %s\n", resp.OK, resp.Message)
```

### Add a New Station

```go
resp, err := rb.AddStation(ctx, baseURL, rb.AddStationRequest{
    Name:        "My Radio Station",
    URL:         "https://stream.example.com/live",
    Homepage:    "https://example.com",
    CountryCode: "US",
    Language:    "english",
    Tags:        "music,pop,rock",
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Added station with UUID: %s\n", resp.UUID)
```

### Check Station Health

```go
// Get check results for a station
checks, err := rb.StationCheckResults(ctx, baseURL, "96202c3a-0601-11e8-ae97-52543be04c81", rb.StationCheckOptions{
    Limit: 5,
})

// Get detailed check steps
steps, err := rb.StationCheckSteps(ctx, baseURL, []string{"96202c3a-0601-11e8-ae97-52543be04c81"})
```

## API Reference

### Station Functions

| Function | Description |
|----------|-------------|
| `Stations` | Get all stations with pagination and ordering |
| `StationsSearch` | Search by kind (name, codec, country, language, tag, etc.) |
| `StationSearch` | Advanced search with multiple filters |
| `StationsByUUID` | Get stations by their UUIDs |
| `StationsByURL` | Find stations by streaming URL |
| `StationsByVotes` | Get top voted stations |
| `StationsByClicks` | Get most clicked stations |
| `StationsByRecentClick` | Get recently played stations |
| `StationsByRecentlyChanged` | Get recently modified stations |
| `StationsChanged` | Track changes to stations |
| `StationsBroken` | Get broken/unavailable stations |

### Metadata Functions

| Function | Description |
|----------|-------------|
| `Countries` | List countries with station counts |
| `Codecs` | List available audio codecs |
| `States` | List states/regions by country |
| `Languages` | List available languages |
| `Tags` | List available tags |
| `StreamingServers` | Get Radio Browser network servers |

### Action Functions

| Function | Description |
|----------|-------------|
| `VoteForStation` | Vote for a station |
| `AddStation` | Add a new station to the database |
| `StationClickCounter` | Register a click/play event |

### Check Functions

| Function | Description |
|----------|-------------|
| `StationCheckResults` | Get availability check results |
| `StationCheckSteps` | Get detailed check steps |
| `StationClicks` | Get click history for a station |

## Project Structure

```
radiobrowser-go/
├── go.mod                 # Module definition (Go 1.22+)
├── types.go               # Type definitions (Station, Country, etc.)
├── stations.go            # Station query functions
├── stations_options.go    # Search options and constants
├── meta.go                # Metadata query functions
├── actions.go             # Action functions (vote, add station)
├── checks.go              # Station health check functions
├── clicks.go              # Click tracking functions
├── *_test.go              # Test files
├── examples/              # Usage examples
│   ├── basic/             # Basic station search
│   ├── search/            # Advanced search
│   ├── metadata/          # Metadata queries
│   └── topstations/       # Popular stations
└── internal/
    └── transport/         # HTTP transport layer
        ├── get.go         # HTTP GET helper
        └── post_form.go   # HTTP POST form helper
```

## Radio Browser API Servers

Radio Browser runs on multiple servers. You can use any of the following base URLs:

- `https://de1.api.radio-browser.info`
- `https://nl1.api.radio-browser.info`
- `https://at1.api.radio-browser.info`

For automatic server discovery, see the [Radio Browser API documentation](https://api.radio-browser.info/).

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

```
MIT License

Copyright (c) 2025 Yuri Sergeev

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
```

## Related Links

- [Radio Browser Website](https://www.radio-browser.info/)
- [Radio Browser API Documentation](https://api.radio-browser.info/)
- [Go Package Documentation](https://pkg.go.dev/github.com/randomtoy/radiobrowser-go)
