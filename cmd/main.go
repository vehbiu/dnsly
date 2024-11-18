package main

import (
	"github.com/vehbiu/dnsly/internal/config"
	"github.com/vehbiu/dnsly/internal/dns"
	"github.com/vehbiu/dnsly/server"
	"log"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dnsServer := dns.NewServer(cfg)
	webServer := server.NewServer(cfg, dnsServer.GetBlockedRequests())

	/* HTTP server is on a GoRoutine as it's not top priority */
	go func() {
		log.Printf("Starting DNS server on :53")
		if err := dnsServer.Start(); err != nil {
			log.Fatalf("Failed to start DNS server: %v", err)
		}
	}()

	/* Web server */
	log.Println("Starting web interface and DoH endpoint on :8080")
	if err := webServer.Start(":8080"); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}
