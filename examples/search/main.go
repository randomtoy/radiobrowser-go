// Example: Advanced station search with multiple filters
//
// This example demonstrates how to use the advanced search API
// to find stations matching multiple criteria simultaneously.
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

	// Search for jazz stations in the US with high bitrate
	fmt.Printf("=== Jazz stations in the US (MP3, 128+ kbps) ===\n\n")

	stations, err := rb.StationSearch(ctx, baseURL, rb.StationSearchOptions{
		CountryCode: "US",
		Tag:         "jazz",
		Codec:       "MP3",
		BitrateMin:  128,
		Limit:       10,
		Order:       rb.StationOrderClickCount,
		Reverse:     true,
		HideBroken:  true,
	})
	if err != nil {
		log.Fatalf("Failed to search stations: %v", err)
	}

	printStations(stations)

	// Search by multiple tags
	fmt.Printf("\n=== Classical music stations with HTTPS ===\n\n")

	isHTTPS := true
	classical, err := rb.StationSearch(ctx, baseURL, rb.StationSearchOptions{
		TagList:    []string{"classical", "orchestra"},
		IsHTTPS:    &isHTTPS,
		Limit:      5,
		Order:      rb.StationOrderVotes,
		Reverse:    true,
		HideBroken: true,
	})
	if err != nil {
		log.Fatalf("Failed to search stations: %v", err)
	}

	printStations(classical)

	// Search by name
	fmt.Printf("\n=== Stations with 'BBC' in name ===\n\n")

	bbc, err := rb.StationsSearch(
		ctx,
		baseURL,
		rb.StationByName,
		"BBC",
		rb.StationsOptions{
			Limit:      5,
			Order:      rb.StationOrderVotes,
			Reverse:    true,
			HideBroken: true,
		},
	)
	if err != nil {
		log.Fatalf("Failed to search stations: %v", err)
	}

	printStations(bbc)

	// Search within 50 km of Amsterdam. Coordinates and radius may be zero.
	lat, lon, radius := 52.37, 4.90, 50000.0
	fmt.Printf("\n=== Stations within 50 km of Amsterdam ===\n\n")
	nearby, err := rb.StationSearch(ctx, baseURL, rb.StationSearchOptions{
		GeoLat: &lat, GeoLong: &lon, GeoDistance: &radius,
		Limit: 5, HideBroken: true,
	})
	if err != nil {
		log.Fatalf("Failed geographic search: %v", err)
	}
	printStations(nearby)
}

func printStations(stations []rb.Station) {
	if len(stations) == 0 {
		fmt.Println("No stations found.")
		return
	}

	for i, s := range stations {
		fmt.Printf("%d. %s (%s)\n", i+1, s.Name, s.CountryCode)
		fmt.Printf("   Codec: %s @ %d kbps | Tags: %s\n", s.Codec, s.Bitrate, truncate(s.Tags, 50))
		fmt.Printf("   Votes: %d | Clicks (24h): %d\n", s.Votes, s.ClickCount)
	}
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
}
