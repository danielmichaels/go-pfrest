package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/danielmichaels/go-pfrest"
)

func main() {
	url := flag.String("url", "", "pfSense base URL (required)")
	apiKey := flag.String("api-key", "", "API key (required)")
	insecure := flag.Bool("insecure", true, "Skip TLS verification")
	flag.Parse()

	if *url == "" || *apiKey == "" {
		flag.Usage()
		os.Exit(1)
	}

	client, err := pfrest.NewClient(pfrest.Config{
		BaseURL:            *url,
		InsecureSkipVerify: *insecure,
		APIKey:             *apiKey,
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	resp, err := client.Raw().GetSystemVersionEndpointWithResponse(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := pfrest.CheckResponse(resp.HTTPResponse); err != nil {
		log.Fatal(err)
	}

	if resp.JSON200 != nil && resp.JSON200.Data != nil {
		data := resp.JSON200.Data
		version := "<unknown>"
		if data.Version != nil {
			version = *data.Version
		}
		fmt.Printf("pfSense version: %s\n", version)
	}
}
