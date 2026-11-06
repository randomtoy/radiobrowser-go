// Example: Basic station search by country
//
// This example demonstrates how to search for radio stations
// by country and display basic information about them.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	rb "github.com/randomtoy/radiobrowser-go"
)

const baseURL = "https://de1.api.radio-browser.info"

func main() {
	ctx, cancel := context.WithTimeout(rb.WithUserAgent(context.Background(), "radiobrowser-go-examples/0.2.0"), 30*time.Second)
	defer cancel()

	// Search for German radio stations, sorted by votes
	stations, err := rb.StationsSearch(
		ctx,
		baseURL,
		rb.StationByCountryCodeExact,
		"DE",
		rb.StationsOptions{
			Limit:      10,
			Order:      rb.StationOrderVotes,
			Reverse:    true, // Descending order (most votes first)
			HideBroken: true,
		},
	)
	if err != nil {
		log.Fatalf("Failed to search stations: %v", err)
	}

	fmt.Printf("Found %d stations in Germany:\n\n", len(stations))

	for i, s := range stations {
		fmt.Printf("%d. %s\n", i+1, s.Name)
		fmt.Printf("   Country: %s\n", s.CountryCode)
		fmt.Printf("   Codec: %s, Bitrate: %d kbps\n", s.Codec, s.Bitrate)
		fmt.Printf("   Votes: %d, Clicks (24h): %d\n", s.Votes, s.ClickCount)
		fmt.Printf("   URL: %s\n", s.URL)
		fmt.Println()
	}
}
