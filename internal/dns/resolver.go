package dns

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/miekg/dns"
)

type Resolver interface {
	Resolve(domain string, qtype uint16) ([]dns.RR, error)
}

type DoHResolver struct {
	endpoint string
}

func NewDoHResolver(endpoint string) *DoHResolver {
	return &DoHResolver{endpoint: endpoint}
}

func (r *DoHResolver) Resolve(domain string, qtype uint16) ([]dns.RR, error) {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), qtype)

	packed, err := msg.Pack()
	if err != nil {
		return nil, fmt.Errorf("failed to pack DNS message: %v", err)
	}

	b64 := base64.RawURLEncoding.EncodeToString(packed)
	req, err := http.NewRequest(http.MethodGet, r.endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/dns-message")
	q := req.URL.Query()
	q.Add("dns", b64)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DoH request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := new(dns.Msg)
	if err := response.Unpack(body); err != nil {
		return nil, err
	}

	return response.Answer, nil
}
