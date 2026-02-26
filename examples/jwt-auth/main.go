package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/danielmichaels/go-pfrest"
	"github.com/danielmichaels/go-pfrest/api"
)

func main() {
	url := flag.String("url", "", "pfSense base URL (required)")
	user := flag.String("user", "admin", "Username")
	pass := flag.String("pass", "", "Password (required)")
	insecure := flag.Bool("insecure", true, "Skip TLS verification")
	flag.Parse()

	if *url == "" || *pass == "" {
		flag.Usage()
		os.Exit(1)
	}

	ctx := context.Background()

	basicClient, err := pfrest.NewClient(pfrest.Config{
		BaseURL:            *url,
		InsecureSkipVerify: *insecure,
		BasicAuth: &pfrest.BasicAuthConfig{
			Username: *user,
			Password: *pass,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	tokenResp, err := basicClient.Raw().PostAuthJWTEndpointWithResponse(ctx, api.PostAuthJWTEndpointJSONRequestBody{})
	if err != nil {
		log.Fatal(err)
	}
	if err := pfrest.CheckResponse(tokenResp.HTTPResponse); err != nil {
		log.Fatal(err)
	}

	if tokenResp.JSON200 == nil || tokenResp.JSON200.Data == nil || tokenResp.JSON200.Data.Token == nil {
		log.Fatal("no JWT token in response")
	}
	token := *tokenResp.JSON200.Data.Token
	fmt.Printf("JWT token obtained (length=%d)\n", len(token))

	jwtClient, err := pfrest.NewClient(pfrest.Config{
		BaseURL:            *url,
		InsecureSkipVerify: *insecure,
		JWTToken:           token,
	})
	if err != nil {
		log.Fatal(err)
	}

	versionResp, err := jwtClient.Raw().GetSystemVersionEndpointWithResponse(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if err := pfrest.CheckResponse(versionResp.HTTPResponse); err != nil {
		log.Fatal(err)
	}

	if versionResp.JSON200 != nil && versionResp.JSON200.Data != nil {
		version := "<unknown>"
		if versionResp.JSON200.Data.Version != nil {
			version = *versionResp.JSON200.Data.Version
		}
		fmt.Printf("pfSense version (via JWT): %s\n", version)
	}
}
