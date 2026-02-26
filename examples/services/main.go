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
	user := flag.String("user", "admin", "Username")
	pass := flag.String("pass", "", "Password")
	apiKey := flag.String("api-key", "", "API key")
	insecure := flag.Bool("insecure", true, "Skip TLS verification")
	flag.Parse()

	if *url == "" {
		flag.Usage()
		os.Exit(1)
	}

	cfg := pfrest.Config{
		BaseURL:            *url,
		InsecureSkipVerify: *insecure,
	}
	switch {
	case *apiKey != "":
		cfg.APIKey = *apiKey
	case *pass != "":
		cfg.BasicAuth = &pfrest.BasicAuthConfig{Username: *user, Password: *pass}
	default:
		log.Fatal("provide -api-key or -pass for authentication")
	}

	client, err := pfrest.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	resp, err := client.Raw().GetStatusServicesEndpointWithResponse(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := pfrest.CheckResponse(resp.HTTPResponse); err != nil {
		log.Fatal(err)
	}

	if resp.JSON200 == nil || resp.JSON200.Data == nil {
		fmt.Println("No services data returned")
		return
	}

	services := *resp.JSON200.Data
	fmt.Printf("Total services: %d\n\n", len(services))

	for _, svc := range services {
		name := ""
		if svc.Name != nil {
			name = *svc.Name
		}
		descr := ""
		if svc.Description != nil {
			descr = *svc.Description
		}
		running := false
		if svc.Status != nil {
			running = *svc.Status
		}
		state := "stopped"
		if running {
			state = "running"
		}
		fmt.Printf("  %-20s %-10s %s\n", name, state, descr)
	}
}
