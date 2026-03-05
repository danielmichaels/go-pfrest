# go-pfrest

> [!WARNING]
> In development; not ready for production use.

Go SDK for the [pfSense REST API](https://github.com/jaredhendrickson13/pfsense-api) v2.

Generated from the pfSense REST API v2.7.2 OpenAPI specification using [Fern](https://buildwithfern.com). Provides modular, per-service clients with built-in retries and typed errors.

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

    pfrest "github.com/danielmichaels/go-pfrest"
    "github.com/danielmichaels/go-pfrest/pkg/client/client"
    "github.com/danielmichaels/go-pfrest/pkg/client/option"
)

func main() {
    c := client.NewClient(
        option.WithBaseURL("https://192.168.1.1"),
        option.WithBasicAuth("admin", "pfsense"),
        option.WithHTTPClient(pfrest.TLSClient(true)),
    )

    resp, err := c.System.GetSystemVersionEndpoint(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    if resp.Data != nil && resp.Data.Version != nil {
        fmt.Printf("pfSense version: %s\n", *resp.Data.Version)
    }
}
```

## Authentication

```go
// Basic Auth
c := client.NewClient(
    option.WithBaseURL("https://192.168.1.1"),
    option.WithBasicAuth("admin", "pfsense"),
)

// API Key
c := client.NewClient(
    option.WithBaseURL("https://192.168.1.1"),
    option.WithAPIKey("your-api-key"),
)

// JWT Bearer Token (obtain via c.Auth.PostAuthJwtEndpoint first)
c := client.NewClient(
    option.WithBaseURL("https://192.168.1.1"),
    option.WithHTTPHeader(http.Header{
        "Authorization": []string{"Bearer " + token},
    }),
)
```

## Usage

The client is organized by service — each pfSense subsystem has its own sub-client:

```go
c.Firewall.GetFirewallRulesEndpoint(ctx, &pfclientapi.GetFirewallRulesEndpointRequest{})
c.Status.GetStatusSystemEndpoint(ctx)
c.System.GetSystemVersionEndpoint(ctx)
c.Diagnostics.GetDiagnosticsArpTableEndpoint(ctx, &pfclientapi.GetDiagnosticsArpTableEndpointRequest{})
c.Services.GetServicesUnboundSettingsEndpoint(ctx)
```

## Error Handling

Errors are returned as typed Go errors. Non-2xx responses are automatically parsed:

```go
resp, err := c.Firewall.GetFirewallRulesEndpoint(ctx, &pfclientapi.GetFirewallRulesEndpointRequest{})
if err != nil {
    // err contains status code and parsed error body
    log.Fatal(err)
}
```

## Retries

Built-in retry with exponential backoff:

```go
c := client.NewClient(
    option.WithBaseURL("https://192.168.1.1"),
    option.WithBasicAuth("admin", "pfsense"),
    option.WithMaxAttempts(3),
)
```

## TLS

pfSense typically uses self-signed certificates. Use the `TLSClient` helper:

```go
c := client.NewClient(
    option.WithBaseURL("https://192.168.1.1"),
    option.WithBasicAuth("admin", "pfsense"),
    option.WithHTTPClient(pfrest.TLSClient(true)), // skip TLS verification
)
```

## Examples

See the [examples/](examples/) directory:

| Example | Description |
|---------|-------------|
| `basic-auth` | Connect with basic auth, list firewall rules |
| `api-key` | Connect with API key, get system version |
| `jwt-auth` | Obtain JWT token, then use it for subsequent calls |
| `firewall` | List firewall rules with type, protocol, interface |
| `services` | List all services with running status |
| `status` | System info, DHCP leases, ARP table |

Run an example:

```bash
go run ./examples/basic-auth -url https://192.168.1.1:10443 -user admin -pass pfsense -insecure
```

## Development

Requires [Fern CLI](https://docs.buildwithfern.com/), Python 3, and [Task](https://taskfile.dev).

```bash
task generate   # Regenerate pkg/client/ (see pipeline below)
task test       # Run tests
task lint       # Run golangci-lint
task build      # Build all packages and examples
task check      # lint + test + build
```

### Code generation pipeline

`task generate` runs these steps in order:

1. **specclean** — `tools/specclean/clean_pfsense_spec.py` normalises the upstream spec and writes `specs/v2.7/openapi-clean.json` (not committed).
2. **fern generate** — reads `openapi-clean.json` plus `specs/v2.7/overlay.yaml` and writes `pkg/client/`.
3. **patch** — `task generate:patch` applies `sed` fixes for known Fern codegen bugs that can't be handled via overlay (e.g. [Basic Auth header format](https://github.com/fern-api/fern/issues/6510)).

Never edit `pkg/client/` by hand — changes will be overwritten on the next `task generate`.

### OpenAPI overlay

`specs/v2.7/openapi.json` is the unmodified upstream pfSense REST API spec. Where the spec diverges from real API behaviour, patches live in `specs/v2.7/overlay.yaml` following the [OpenAPI Overlay Specification](https://spec.openapis.org/overlay/v1.0.0). Fern applies this overlay automatically during step 2 above.

To add a new patch, append an action to `overlay.yaml` with a JSONPath `target` and either an `update` (merge) or `remove: true` operation, with a comment explaining the discrepancy.

## License

Apache 2.0 — see [LICENSE](LICENSE).
