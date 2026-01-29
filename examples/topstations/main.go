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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	opts := rb.StationsOptions{
		Limit:      10,
		HideBroken: true,
	}

	// Top voted stations
	fmt.Println("=== Top 10 Stations by Votes ===\n")

	topVoted, err := rb.StationsByVotes(ctx, baseURL, opts)
	if err != nil {
		log.Fatalf("Failed to get top voted stations: %v", err)
	}

	for i, s := range topVoted {
		fmt.Printf("%2d. %-35s | %6d votes | %s\n",
			i+1, truncate(s.Name, 35), s.Votes, s.Country)
	}

	// Most clicked stations (all time)
	fmt.Println("\n=== Top 10 Stations by Total Clicks ===\n")

	topClicked, err := rb.StationsByClicks(ctx, baseURL, opts)
	if err != nil {
		log.Fatalf("Failed to get top clicked stations: %v", err)
	}

	for i, s := range topClicked {
		fmt.Printf("%2d. %-35s | %6d clicks | %s\n",
			i+1, truncate(s.Name, 35), s.ClickCount, s.Country)
	}

	// Recently played stations
	fmt.Println("\n=== Recently Played Stations ===\n")

	recentlyPlayed, err := rb.StationsByRecentClick(ctx, baseURL, opts)
	if err != nil {
		log.Fatalf("Failed to get recently played stations: %v", err)
	}

	for i, s := range recentlyPlayed {
		fmt.Printf("%2d. %-35s | %s | %s\n",
			i+1, truncate(s.Name, 35), s.ClickTimestamp, s.Country)
	}

	// Trending stations (by click trend)
	fmt.Println("\n=== Trending Stations (by Click Trend) ===\n")

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
			i+1, truncate(s.Name, 35), s.ClickTrend, s.Country)
	}

	// Recently changed/added stations
	fmt.Println("\n=== Recently Added/Changed Stations ===\n")

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
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
