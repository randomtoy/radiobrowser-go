// Example: Reading Radio Browser API server information.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	rb "github.com/randomtoy/radiobrowser-go"
)

func main() {
	baseURL := flag.String("server", "https://de1.api.radio-browser.info", "Radio Browser API base URL")
	mirrors := flag.Bool("mirrors", false, "demonstrate HTTP mirror lookup (fallback when DNS is unavailable)")
	flag.Parse()
	ctx, cancel := context.WithTimeout(rb.WithUserAgent(context.Background(), "radiobrowser-go-examples/0.2.0"), 15*time.Second)
	defer cancel()

	stats, err := rb.Stats(ctx, *baseURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("API %s: %s; %d stations, %d broken; %d clicks in the last day\n",
		stats.SoftwareVersion, stats.Status, stats.Stations, stats.StationsBroken, stats.ClicksLastDay)

	config, err := rb.Config(ctx, *baseURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Location: %s (%s); cache: %s, TTL: %ds\n",
		config.ServerLocation, config.ServerCountryCode, config.CacheType, config.CacheTTL)

	if *mirrors {
		servers, err := rb.Servers(ctx, *baseURL)
		if err != nil {
			log.Fatal(err)
		}
		for _, server := range servers {
			fmt.Printf("Mirror: %s (%s)\n", server.Name, server.IP)
		}
	}
}
