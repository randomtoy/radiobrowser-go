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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get countries with most stations
	fmt.Println("=== Top 10 Countries by Station Count ===\n")

	countries, err := rb.Countries(ctx, baseURL, "", rb.CountriesOptions{
		Order:   "stationcount",
		Reverse: true,
	})
	if err != nil {
		log.Fatalf("Failed to get countries: %v", err)
	}

	for i, c := range countries[:min(10, len(countries))] {
		fmt.Printf("%2d. %-30s (%s) - %d stations\n", i+1, c.Name, c.ISO3166_1, c.StationCount)
	}

	// Get popular tags
	fmt.Println("\n=== Top 15 Popular Tags ===\n")

	tags, err := rb.Tags(ctx, baseURL, "", rb.TagsOptions{
		Order:   "stationcount",
		Reverse: true,
	})
	if err != nil {
		log.Fatalf("Failed to get tags: %v", err)
	}

	for i, t := range tags[:min(15, len(tags))] {
		fmt.Printf("%2d. %-20s - %d stations\n", i+1, t.Name, t.StationCount)
	}

	// Get available codecs
	fmt.Println("\n=== Available Codecs ===\n")

	codecs, err := rb.Codecs(ctx, baseURL, "", rb.CodecsOptions{
		Order:   "stationcount",
		Reverse: true,
	})
	if err != nil {
		log.Fatalf("Failed to get codecs: %v", err)
	}

	for i, c := range codecs[:min(10, len(codecs))] {
		fmt.Printf("%2d. %-15s - %d stations\n", i+1, c.Name, c.StationCount)
	}

	// Get languages
	fmt.Println("\n=== Top 10 Languages ===\n")

	languages, err := rb.Languages(ctx, baseURL, "", rb.LanguagesOptions{
		Order:   "stationcount",
		Reverse: true,
	})
	if err != nil {
		log.Fatalf("Failed to get languages: %v", err)
	}

	for i, l := range languages[:min(10, len(languages))] {
		fmt.Printf("%2d. %-20s - %d stations\n", i+1, l.Name, l.StationCount)
	}

	// Get states for a specific country
	fmt.Println("\n=== US States with Most Stations ===\n")

	states, err := rb.States(ctx, baseURL, "United States", "", rb.StatesOptions{
		Order:   "stationcount",
		Reverse: true,
	})
	if err != nil {
		log.Fatalf("Failed to get states: %v", err)
	}

	for i, s := range states[:min(10, len(states))] {
		fmt.Printf("%2d. %-25s - %d stations\n", i+1, s.Name, s.StationCount)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
