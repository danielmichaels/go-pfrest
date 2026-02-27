package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

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

	resp, err := c.Firewall.GetFirewallRulesEndpoint(ctx, &pfclientapi.GetFirewallRulesEndpointRequest{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total firewall rules: %d\n\n", len(resp.Data))

	for i, rule := range resp.Data {
		iface := "<none>"
		if rule.Interface != nil {
			iface = strings.Join(rule.Interface, ",")
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
