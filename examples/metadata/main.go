// Example: Fetching metadata (countries, languages, tags, codecs)
//
// This example demonstrates how to retrieve metadata from the
// Radio Browser API, such as available countries, languages,
// popular tags, and supported codecs.
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

	// Get countries with most stations
	fmt.Printf("=== Top 10 Countries by Station Count ===\n\n")

	countries, err := rb.Countries(ctx, baseURL, "", rb.CountriesOptions{
		Order:   "stationcount",
		Reverse: true,
		Limit:   10,
	})
	if err != nil {
		log.Fatalf("Failed to get countries: %v", err)
	}

	for i, c := range countries {
		fmt.Printf("%2d. %-30s (%s) - %d stations\n", i+1, c.Name, c.ISO3166_1, c.StationCount)
	}

	// Get popular tags
	fmt.Printf("\n=== Top 15 Popular Tags ===\n\n")

	tags, err := rb.Tags(ctx, baseURL, "", rb.TagsOptions{
		Order:   "stationcount",
		Reverse: true,
		Limit:   15,
	})
	if err != nil {
		log.Fatalf("Failed to get tags: %v", err)
	}

	for i, t := range tags {
		fmt.Printf("%2d. %-20s - %d stations\n", i+1, t.Name, t.StationCount)
	}

	// Get available codecs
	fmt.Printf("\n=== Available Codecs ===\n\n")

	codecs, err := rb.Codecs(ctx, baseURL, "", rb.CodecsOptions{
		Order:   "stationcount",
		Reverse: true,
		Limit:   10,
	})
	if err != nil {
		log.Fatalf("Failed to get codecs: %v", err)
	}

	for i, c := range codecs {
		fmt.Printf("%2d. %-15s - %d stations\n", i+1, c.Name, c.StationCount)
	}

	// Get languages
	fmt.Printf("\n=== Top 10 Languages ===\n\n")

	languages, err := rb.Languages(ctx, baseURL, "", rb.LanguagesOptions{
		Order:   "stationcount",
		Reverse: true,
		Limit:   10,
	})
	if err != nil {
		log.Fatalf("Failed to get languages: %v", err)
	}

	for i, l := range languages {
		fmt.Printf("%2d. %-20s - %d stations\n", i+1, l.Name, l.StationCount)
	}

	// Get states for a specific country
	fmt.Printf("\n=== US States with Most Stations ===\n\n")

	states, err := rb.States(ctx, baseURL, "The United States Of America", "", rb.StatesOptions{
		Order:   "stationcount",
		Reverse: true,
		Limit:   10,
	})
	if err != nil {
		log.Fatalf("Failed to get states: %v", err)
	}

	for i, s := range states {
		fmt.Printf("%2d. %-25s - %d stations\n", i+1, s.Name, s.StationCount)
	}
}
