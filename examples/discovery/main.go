// Example: Discovering API servers through DNS before the first HTTP request.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	rb "github.com/randomtoy/radiobrowser-go"
)

func main() {
	ctx, cancel := context.WithTimeout(rb.WithUserAgent(context.Background(), "radiobrowser-go-discovery-example/1.0"), 30*time.Second)
	defer cancel()
	servers, err := rb.DiscoverServers(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Discovered %d API servers: %v\n", len(servers), servers)

	for _, baseURL := range servers {
		// Bound each attempt so an unavailable mirror cannot consume the whole deadline.
		attempt, stop := context.WithTimeout(ctx, 5*time.Second)
		stats, err := rb.Stats(attempt, baseURL)
		stop()
		if err != nil {
			log.Printf("%s: %v", baseURL, err)
			continue
		}
		fmt.Printf("%s: API %s; %d stations\n", baseURL, stats.SoftwareVersion, stats.Stations)
		return
	}
	log.Fatal("no discovered API server responded successfully")
}
