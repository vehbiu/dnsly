package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type Config struct {
	BlockedDomains []string `json:"blocked_domains"`
	SearchDNS      string   `json:"search_dns"`
	DoHEndpoint    string   `json:"doh_endpoint"`
	HTTPPassword   string   `json:"http_password"`
}

func Load(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	if config.SearchDNS == "" {
		config.SearchDNS = "8.8.8.8"
	}
	if config.DoHEndpoint == "" {
		config.DoHEndpoint = "https://dns.google/dns-query"
	}
	if config.HTTPPassword == "" {
		config.HTTPPassword = "admin"
		for i := 0; i < 10; i++ {
			log.Println("[!!!] Warning: HTTP password is set to default value")
		}
	}

	return config, nil
}

func (c *Config) IsDomainBlocked(domain string) bool {
	domain = strings.TrimSuffix(domain, ".")
	for _, blockedDomain := range c.BlockedDomains {
		if strings.TrimSuffix(blockedDomain, ".") == domain {
			return true
		}
	}
	return false
}
