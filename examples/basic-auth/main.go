package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/danielmichaels/go-pfrest"
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

	client, err := pfrest.NewClient(pfrest.Config{
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

	ctx := context.Background()
	resp, err := client.Raw().GetFirewallRulesEndpointWithResponse(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := pfrest.CheckResponse(resp.HTTPResponse); err != nil {
		log.Fatal(err)
	}

	if resp.JSON200 != nil && resp.JSON200.Data != nil {
		fmt.Printf("Firewall rules: %d\n", len(*resp.JSON200.Data))
		for i, rule := range *resp.JSON200.Data {
			iface := "<none>"
			if rule.Interface != nil {
				iface = strings.Join(*rule.Interface, ",")
			}
			descr := ""
			if rule.Descr != nil {
				descr = *rule.Descr
			}
			fmt.Printf("  [%d] interface=%s descr=%s\n", i, iface, descr)
		}
	}
}
