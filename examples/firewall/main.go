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

	resp, err := client.Raw().GetFirewallRulesEndpointWithResponse(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := pfrest.CheckResponse(resp.HTTPResponse); err != nil {
		log.Fatal(err)
	}

	if resp.JSON200 == nil || resp.JSON200.Data == nil {
		fmt.Println("No firewall rules data returned")
		return
	}

	rules := *resp.JSON200.Data
	fmt.Printf("Total firewall rules: %d\n\n", len(rules))

	for i, rule := range rules {
		iface := "<none>"
		if rule.Interface != nil {
			iface = strings.Join(*rule.Interface, ",")
		}
		descr := ""
		if rule.Descr != nil {
			descr = *rule.Descr
		}
		ruleType := ""
		if rule.Type != nil {
			ruleType = string(*rule.Type)
		}
		proto := ""
		if rule.Ipprotocol != nil {
			proto = string(*rule.Ipprotocol)
		}
		fmt.Printf("  [%d] type=%-5s proto=%-5s iface=%-10s descr=%s\n",
			i, ruleType, proto, iface, descr)
	}
}
