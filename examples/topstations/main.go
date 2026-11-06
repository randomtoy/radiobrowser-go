// Example: Getting popular stations
//
// This example demonstrates how to retrieve the most popular
// radio stations based on different metrics: votes, clicks,
// recent activity, and trending.
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

	opts := rb.StationsOptions{
		Limit:      10,
		HideBroken: true,
	}

	// Top voted stations
	fmt.Printf("=== Top 10 Stations by Votes ===\n\n")

	topVoted, err := rb.StationsByVotes(ctx, baseURL, opts)
	if err != nil {
		log.Fatalf("Failed to get top voted stations: %v", err)
	}

	for i, s := range topVoted {
		fmt.Printf("%2d. %-35s | %6d votes | %s\n",
			i+1, truncate(s.Name, 35), s.Votes, s.CountryCode)
	}

	// Most clicked stations in the last 24 hours
	fmt.Printf("\n=== Top 10 Stations by Clicks in the Last 24 Hours ===\n\n")

	topClicked, err := rb.StationsByClicks(ctx, baseURL, opts)
	if err != nil {
		log.Fatalf("Failed to get top clicked stations: %v", err)
	}

	for i, s := range topClicked {
		fmt.Printf("%2d. %-35s | %6d clicks | %s\n",
			i+1, truncate(s.Name, 35), s.ClickCount, s.CountryCode)
	}

	// Recently played stations
	fmt.Printf("\n=== Recently Played Stations ===\n\n")

	recentlyPlayed, err := rb.StationsByRecentClick(ctx, baseURL, opts)
	if err != nil {
		log.Fatalf("Failed to get recently played stations: %v", err)
	}

	for i, s := range recentlyPlayed {
		fmt.Printf("%2d. %-35s | %s | %s\n",
			i+1, truncate(s.Name, 35), s.ClickTimestamp, s.CountryCode)
	}

	// Trending stations (by click trend)
	fmt.Printf("\n=== Trending Stations (by Click Trend) ===\n\n")

	trending, err := rb.Stations(ctx, baseURL, rb.StationsOptions{
		Limit:      10,
		Order:      rb.StationOrderClickTrend,
		Reverse:    true,
		HideBroken: true,
	})
	if err != nil {
		log.Fatalf("Failed to get trending stations: %v", err)
	}

	for i, s := range trending {
		fmt.Printf("%2d. %-35s | trend: %+d | %s\n",
			i+1, truncate(s.Name, 35), s.ClickTrend, s.CountryCode)
	}

	// Recently changed/added stations
	fmt.Printf("\n=== Recently Added/Changed Stations ===\n\n")

	recentlyChanged, err := rb.StationsByRecentlyChanged(ctx, baseURL, opts)
	if err != nil {
		log.Fatalf("Failed to get recently changed stations: %v", err)
	}

	for i, s := range recentlyChanged {
		fmt.Printf("%2d. %-35s | %s\n",
			i+1, truncate(s.Name, 35), s.LastChangeTime)
	}
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
}
