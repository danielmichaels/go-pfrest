// Package pfrest provides helpers for the pfSense REST API Go client.
//
// The generated client lives in the pkg/client sub-module. Import it directly:
//
//	import (
//		"github.com/danielmichaels/go-pfrest/pkg/client/client"
//		"github.com/danielmichaels/go-pfrest/pkg/client/option"
//	)
//
//	c := client.NewClient(
//		option.WithBaseURL("https://pfsense.local"),
//		option.WithBasicAuth("admin", "password"),
//		option.WithHTTPClient(pfrest.TLSClient(true)),
//	)
//
//	rules, err := c.Firewall.GetFirewallRulesEndpoint(ctx, &pfclientapi.GetFirewallRulesEndpointRequest{})
//
// This root package provides [TLSClient] for the common pfSense self-signed
// certificate scenario. Pass the result to option.WithHTTPClient.
package pfrest
