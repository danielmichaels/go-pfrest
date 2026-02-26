package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

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

	fmt.Println("=== System Status ===")
	sysResp, err := client.Raw().GetStatusSystemEndpointWithResponse(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if err := pfrest.CheckResponse(sysResp.HTTPResponse); err != nil {
		log.Fatal(err)
	}
	if sysResp.JSON200 != nil && sysResp.JSON200.Data != nil {
		s := sysResp.JSON200.Data
		if s.Platform != nil {
			fmt.Printf("  Platform:  %s\n", *s.Platform)
		}
		if s.Uptime != nil {
			fmt.Printf("  Uptime:    %s\n", *s.Uptime)
		}
		if s.CPUModel != nil {
			fmt.Printf("  CPU:       %s\n", *s.CPUModel)
		}
		if s.CPUUsage != nil {
			fmt.Printf("  CPU usage: %.1f%%\n", *s.CPUUsage)
		}
		if s.MemUsage != nil {
			fmt.Printf("  Mem usage: %.1f%%\n", *s.MemUsage)
		}
		if s.DiskUsage != nil {
			fmt.Printf("  Disk:      %.1f%%\n", *s.DiskUsage)
		}
		if s.TempC != nil {
			fmt.Printf("  Temp:      %.1f°C\n", *s.TempC)
		}
	}

	fmt.Println("\n=== DHCP Leases ===")
	dhcpResp, err := client.Raw().GetStatusDHCPServerLeasesEndpointWithResponse(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := pfrest.CheckResponse(dhcpResp.HTTPResponse); err != nil {
		log.Fatal(err)
	}
	if dhcpResp.JSON200 != nil && dhcpResp.JSON200.Data != nil {
		leases := *dhcpResp.JSON200.Data
		fmt.Printf("  Total leases: %d\n", len(leases))
		for _, l := range leases {
			ip := ""
			if l.IP != nil {
				ip = *l.IP
			}
			mac := ""
			if l.Mac != nil {
				mac = *l.Mac
			}
			hostname := ""
			if l.Hostname != nil {
				hostname = *l.Hostname
			}
			fmt.Printf("  %-16s %-18s %s\n", ip, mac, hostname)
		}
	}

	fmt.Println("\n=== ARP Table ===")
	arpResp, err := client.Raw().GetDiagnosticsARPTableEndpointWithResponse(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := pfrest.CheckResponse(arpResp.HTTPResponse); err != nil {
		log.Fatal(err)
	}
	if arpResp.JSON200 != nil && arpResp.JSON200.Data != nil {
		entries := *arpResp.JSON200.Data
		fmt.Printf("  Total entries: %d\n", len(entries))
		for _, e := range entries {
			ip := ""
			if e.IPAddress != nil {
				ip = *e.IPAddress
			}
			mac := ""
			if e.MacAddress != nil {
				mac = *e.MacAddress
			}
			iface := ""
			if e.Interface != nil {
				iface = *e.Interface
			}
			fmt.Printf("  %-16s %-18s %s\n", ip, mac, iface)
		}
	}
}
