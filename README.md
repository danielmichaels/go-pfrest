# go-pfrest

Go SDK for the [pfSense REST API](https://github.com/jaredhendrickson13/pfsense-api) v2.

Generated from the pfSense REST API v2.7.2 OpenAPI specification covering all 258 endpoints (677 operations) with fully typed request/response structs.

## Installation

```bash
go get github.com/danielmichaels/go-pfrest
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/danielmichaels/go-pfrest"
)

func main() {
    client, err := pfrest.NewClient(pfrest.Config{
        BaseURL:            "192.168.1.1",
        InsecureSkipVerify: true,
        BasicAuth: &pfrest.BasicAuthConfig{
            Username: "admin",
            Password: "pfsense",
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    resp, err := client.Raw().GetSystemVersionEndpointWithResponse(ctx)
    if err != nil {
        log.Fatal(err)
    }

    if resp.JSON200 != nil && resp.JSON200.Data != nil {
        fmt.Printf("pfSense version: %s\n", *resp.JSON200.Data.Version)
    }
}
```

## Authentication

Three methods are supported. Provide exactly one:

```go
// Basic Auth
cfg := pfrest.Config{
    BaseURL: "192.168.1.1",
    BasicAuth: &pfrest.BasicAuthConfig{
        Username: "admin",
        Password: "pfsense",
    },
}

// API Key
cfg := pfrest.Config{
    BaseURL: "192.168.1.1",
    APIKey:  "your-api-key",
}

// JWT Bearer Token
cfg := pfrest.Config{
    BaseURL:  "192.168.1.1",
    JWTToken: "your-jwt-token",
}
```

The `BaseURL` automatically prepends `https://` if no scheme is provided.

## Usage

All 677 generated operations are accessible via `client.Raw()`, which returns the oapi-codegen `ClientWithResponses`:

```go
// List firewall rules
resp, _ := client.Raw().GetFirewallRulesEndpointWithResponse(ctx, nil)

// Get system status
resp, _ := client.Raw().GetStatusSystemEndpointWithResponse(ctx)

// Restart a service
client.Raw().PostStatusServiceEndpointWithResponse(ctx, api.PostStatusServiceEndpointJSONRequestBody{
    Name:   ptr("unbound"),
    Action: ptr(api.PostStatusServiceEndpointJSONBodyActionRestart),
})
```

## Error Handling

Two approaches:

**Typed responses** — check each status code field directly:

```go
resp, err := client.Raw().GetFirewallRulesEndpointWithResponse(ctx, nil)
if resp.JSON200 != nil {
    // success
}
if resp.JSON400 != nil {
    // validation error with typed fields
    fmt.Println(resp.JSON400.Message)
}
```

**Convenience helper** — parse any non-2xx into `*pfrest.APIError`:

```go
if err := pfrest.CheckResponse(resp.HTTPResponse); err != nil {
    var apiErr *pfrest.APIError
    if errors.As(err, &apiErr) {
        fmt.Println(apiErr.ResponseID) // e.g. "FIREWALL_ALIAS_NAME_EXISTS"
    }
}
```

## TLS

pfSense typically uses self-signed certificates. Set `InsecureSkipVerify: true` or provide a custom `*http.Client`:

```go
cfg := pfrest.Config{
    BaseURL:            "192.168.1.1",
    InsecureSkipVerify: true,
    // OR
    HTTPClient: yourCustomClient,
}
```

## Examples

See the [examples/](examples/) directory:

| Example | Description |
|---------|-------------|
| `basic-auth` | Connect with basic auth, list firewall rules |
| `api-key` | Connect with API key, get system version |
| `jwt-auth` | Obtain JWT token, then use it for subsequent calls |
| `firewall` | List firewall rules with details |
| `services` | List all services with status |
| `status` | System info, DHCP leases, ARP table |

Run an example:

```bash
go run ./examples/basic-auth -url 192.168.1.1:10443 -user admin -pass pfsense
```

## Development

Requires [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) v2.5.1 and [Task](https://taskfile.dev).

```bash
# Regenerate from OpenAPI spec
task generate

# Run tests
task test

# Lint
task lint

# Build everything
task build
```

The spec preprocessor (`tools/specprep`) simplifies the raw 11MB OpenAPI spec to eliminate allOf/oneOf compositions in error responses that cause oapi-codegen to hang.

## License

Apache 2.0 — see [LICENSE](LICENSE).
