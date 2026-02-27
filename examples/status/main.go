package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

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

	fmt.Println("=== System Status ===")
	sysResp, err := c.Status.GetStatusSystemEndpoint(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if sysResp.Data != nil {
		s := sysResp.Data
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
	dhcpResp, err := c.Status.GetStatusDhcpServerLeasesEndpoint(ctx, &pfclientapi.GetStatusDhcpServerLeasesEndpointRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Total leases: %d\n", len(dhcpResp.Data))
	for _, l := range dhcpResp.Data {
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

	fmt.Println("\n=== ARP Table ===")
	arpResp, err := c.Diagnostics.GetDiagnosticsArpTableEndpoint(ctx, &pfclientapi.GetDiagnosticsArpTableEndpointRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Total entries: %d\n", len(arpResp.Data))
	for _, e := range arpResp.Data {
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
