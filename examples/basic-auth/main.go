package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	pfrest "github.com/danielmichaels/go-pfrest"
	client "github.com/danielmichaels/go-pfrest/pkg/client/client"
	pfclientapi "github.com/danielmichaels/go-pfrest/pkg/client"
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

	c := client.NewClient(
		option.WithBaseURL(*url),
		option.WithBasicAuth(*user, *pass),
		option.WithHTTPClient(pfrest.TLSClient(*insecure)),
	)

	ctx := context.Background()
	resp, err := c.Firewall.GetFirewallRulesEndpoint(ctx, &pfclientapi.GetFirewallRulesEndpointRequest{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Firewall rules: %d\n", len(resp.Data))
	for i, rule := range resp.Data {
		iface := "<none>"
		if rule.Interface != nil {
			iface = strings.Join(rule.Interface, ",")
		}
		descr := ""
		if rule.Descr != nil {
			descr = *rule.Descr
		}
		fmt.Printf("  [%d] interface=%s descr=%s\n", i, iface, descr)
	}
}
