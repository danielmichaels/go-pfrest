package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	pfrest "github.com/danielmichaels/go-pfrest"
	pfclientapi "github.com/danielmichaels/go-pfrest/pkg/client"
	client "github.com/danielmichaels/go-pfrest/pkg/client/client"
	"github.com/danielmichaels/go-pfrest/pkg/client/option"
)

func main() {
	url := flag.String("url", "", "pfSense base URL (required)")
	user := flag.String("user", "admin", "Username")
	pass := flag.String("pass", "", "Password")
	apiKey := flag.String("api-key", "", "API key")
	insecure := flag.Bool("insecure", true, "Skip TLS verification")
	flag.Parse()

	if *url == "" {
		flag.Usage()
		os.Exit(1)
	}

	opts := []option.RequestOption{
		option.WithBaseURL(*url),
		option.WithHTTPClient(pfrest.TLSClient(*insecure)),
	}
	switch {
	case *apiKey != "":
		opts = append(opts, option.WithAPIKey(*apiKey))
	case *pass != "":
		opts = append(opts, option.WithBasicAuth(*user, *pass))
	default:
		log.Fatal("provide -api-key or -pass for authentication")
	}

	c := client.NewClient(opts...)
	ctx := context.Background()

	resp, err := c.Status.GetStatusServicesEndpoint(ctx, &pfclientapi.GetStatusServicesEndpointRequest{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total services: %d\n\n", len(resp.Data))

	for _, svc := range resp.Data {
		name := ""
		if svc.Name != nil {
			name = *svc.Name
		}
		descr := ""
		if svc.Description != nil {
			descr = *svc.Description
		}
		state := "stopped"
		if svc.Status != nil && *svc.Status {
			state = "running"
		}
		fmt.Printf("  %-20s %-10s %s\n", name, state, descr)
	}
}
