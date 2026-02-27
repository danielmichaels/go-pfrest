package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	pfrest "github.com/danielmichaels/go-pfrest"
	client "github.com/danielmichaels/go-pfrest/pkg/client/client"
	"github.com/danielmichaels/go-pfrest/pkg/client/option"
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

	c := client.NewClient(
		option.WithBaseURL(*url),
		option.WithAPIKey(*apiKey),
		option.WithHTTPClient(pfrest.TLSClient(*insecure)),
	)

	ctx := context.Background()
	resp, err := c.System.GetSystemVersionEndpoint(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if resp.Data != nil {
		version := "<unknown>"
		if resp.Data.Version != nil {
			version = *resp.Data.Version
		}
		fmt.Printf("pfSense version: %s\n", version)
	}
}
