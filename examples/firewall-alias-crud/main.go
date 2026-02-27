package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	pfrest "github.com/danielmichaels/go-pfrest"
	pfapi "github.com/danielmichaels/go-pfrest/pkg/client"
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

	// 1. List aliases (baseline)
	fmt.Println("==> Listing aliases (baseline)")
	listResp, err := c.Firewall.GetFirewallAliasesEndpoint(ctx, &pfapi.GetFirewallAliasesEndpointRequest{})
	if err != nil {
		log.Fatalf("list: %v", err)
	}
	fmt.Printf("    Found %d aliases\n", len(listResp.Data))

	// 2. Create alias
	fmt.Println("==> Creating alias 'e2e_test_alias'")
	name := "e2e_test_alias"
	aliasType := pfapi.FirewallAliasTypeHost
	descr := "created by e2e example"
	createResp, err := c.Firewall.PostFirewallAliasEndpoint(ctx, &pfapi.PostFirewallAliasEndpointRequest{
		Name:    &name,
		Type:    &aliasType,
		Descr:   &descr,
		Address: []string{"10.99.99.1"},
		Detail:  []string{"e2e entry"},
	})
	if err != nil {
		log.Fatalf("create: %v", err)
	}
	if createResp.Data == nil || createResp.Data.ID == nil {
		log.Fatal("create: no ID in response")
	}
	id := *createResp.Data.ID
	fmt.Printf("    Created: id=%d name=%s\n", id, str(createResp.Data.Name))

	// 3. Get alias by ID
	fmt.Printf("==> Getting alias id=%d\n", id)
	idStr := fmt.Sprintf("%d", id)
	getResp, err := c.Firewall.GetFirewallAliasEndpoint(ctx, &pfapi.GetFirewallAliasEndpointRequest{
		ID: &idStr,
	})
	if err != nil {
		log.Fatalf("get: %v", err)
	}
	fmt.Printf("    Got: name=%s descr=%s addresses=%v\n",
		str(getResp.Data.Name), str(getResp.Data.Descr), getResp.Data.Address)

	// 4. Update alias (patch)
	fmt.Printf("==> Updating alias id=%d\n", id)
	newDescr := "updated by e2e example"
	patchResp, err := c.Firewall.PatchFirewallAliasEndpoint(ctx, &pfapi.PatchFirewallAliasEndpointRequest{
		ID:      id,
		Descr:   &newDescr,
		Address: []string{"10.99.99.1", "10.99.99.2"},
		Detail:  []string{"entry 1", "entry 2"},
	})
	if err != nil {
		log.Fatalf("update: %v", err)
	}
	fmt.Printf("    Updated: descr=%s addresses=%v\n",
		str(patchResp.Data.Descr), patchResp.Data.Address)

	// 5. Get again to verify update
	fmt.Printf("==> Verifying update id=%d\n", id)
	getResp2, err := c.Firewall.GetFirewallAliasEndpoint(ctx, &pfapi.GetFirewallAliasEndpointRequest{
		ID: &idStr,
	})
	if err != nil {
		log.Fatalf("get after update: %v", err)
	}
	fmt.Printf("    Verified: descr=%s addresses=%v\n",
		str(getResp2.Data.Descr), getResp2.Data.Address)

	// 6. Delete alias
	fmt.Printf("==> Deleting alias id=%d\n", id)
	apply := true
	_, err = c.Firewall.DeleteFirewallAliasEndpoint(ctx, &pfapi.DeleteFirewallAliasEndpointRequest{
		ID:    &idStr,
		Apply: &apply,
	})
	if err != nil {
		log.Fatalf("delete: %v", err)
	}
	fmt.Println("    Deleted")

	// 7. List to confirm deletion
	fmt.Println("==> Listing aliases (after delete)")
	listResp2, err := c.Firewall.GetFirewallAliasesEndpoint(ctx, &pfapi.GetFirewallAliasesEndpointRequest{})
	if err != nil {
		log.Fatalf("list after delete: %v", err)
	}
	for _, item := range listResp2.Data {
		if item.Name != nil && *item.Name == name {
			log.Fatal("alias still exists after delete!")
		}
	}
	fmt.Printf("    Confirmed: %d aliases, e2e_test_alias gone\n", len(listResp2.Data))
	fmt.Println("==> PASS: full CRUD cycle complete")
}

func str(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}
