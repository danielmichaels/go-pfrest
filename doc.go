// Package pfrest provides a Go client for the pfSense REST API v2.
//
// The client is generated from the pfSense REST API v2.7.2 OpenAPI specification
// using oapi-codegen and covers all 258 endpoints (677 operations).
//
// Three authentication methods are supported: HTTP Basic Auth, API Key, and JWT
// Bearer Token.
//
// Quick start:
//
//	client, err := pfrest.NewClient(pfrest.Config{
//	    BaseURL:            "https://pfsense.local",
//	    InsecureSkipVerify: true,
//	    BasicAuth: &pfrest.BasicAuthConfig{
//	        Username: "admin",
//	        Password: "pfsense",
//	    },
//	})
//
// Access all generated operations via the Raw() method:
//
//	resp, err := client.Raw().GetFirewallRulesEndpointWithResponse(ctx)
//
// Error handling supports both typed responses and convenience checking:
//
//	if resp.JSON200 != nil {
//	    // Success — use typed data
//	}
//	if resp.JSON400 != nil {
//	    // Validation error — typed error response
//	}
package pfrest
