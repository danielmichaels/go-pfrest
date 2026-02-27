package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	pfrest "github.com/danielmichaels/go-pfrest"
	pfclientapi "github.com/danielmichaels/go-pfrest/pkg/client"
	client "github.com/danielmichaels/go-pfrest/pkg/client/client"
	"github.com/danielmichaels/go-pfrest/pkg/client/option"
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
	httpClient := pfrest.TLSClient(*insecure)

	basicClient := client.NewClient(
		option.WithBaseURL(*url),
		option.WithBasicAuth(*user, *pass),
		option.WithHTTPClient(httpClient),
	)

	tokenResp, err := basicClient.Auth.PostAuthJwtEndpoint(ctx, &pfclientapi.PostAuthJwtEndpointRequest{})
	if err != nil {
		log.Fatal(err)
	}
	if tokenResp.Data == nil || tokenResp.Data.Token == nil {
		log.Fatal("no JWT token in response")
	}
	token := *tokenResp.Data.Token
	fmt.Printf("JWT token obtained (length=%d)\n", len(token))

	jwtClient := client.NewClient(
		option.WithBaseURL(*url),
		option.WithHTTPClient(httpClient),
		option.WithHTTPHeader(http.Header{
			"Authorization": []string{"Bearer " + token},
		}),
	)

	versionResp, err := jwtClient.System.GetSystemVersionEndpoint(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if versionResp.Data != nil {
		version := "<unknown>"
		if versionResp.Data.Version != nil {
			version = *versionResp.Data.Version
		}
		fmt.Printf("pfSense version (via JWT): %s\n", version)
	}
}
