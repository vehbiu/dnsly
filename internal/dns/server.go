package dns

import (
	"log"
	"net"
	"sync"

	"github.com/miekg/dns"
	"github.com/vehbiu/dnsly/internal/config"
	"github.com/vehbiu/dnsly/internal/models"
)

type Server struct {
	cfg             *config.Config
	resolver        Resolver
	blockedRequests []models.BlockedRequest
	mu              sync.RWMutex
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg:      cfg,
		resolver: NewDoHResolver(cfg.DoHEndpoint),
	}
}

func (s *Server) handleDNS(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	for _, q := range r.Question {
		domain := q.Name
		log.Printf("DNS request: %s from %s", domain, w.RemoteAddr().String())

		if s.cfg.IsDomainBlocked(domain) {
			s.handleBlockedDomain(m, q, w.RemoteAddr().String())
			continue
		}

		answer, err := s.resolver.Resolve(domain, q.Qtype)
		if err != nil {
			log.Printf("Failed to resolve %s: %v", domain, err)
			m.SetRcode(r, dns.RcodeServerFailure)
			continue
		}
		m.Answer = append(m.Answer, answer...)
	}

	if err := w.WriteMsg(m); err != nil {
		log.Printf("Failed to write DNS response: %v", err)
	}
}

func (s *Server) handleBlockedDomain(m *dns.Msg, q dns.Question, sourceAddr string) {
	log.Printf("Blocking request for %s", q.Name)

	blockedIP := net.ParseIP("0.0.0.0")
	rr := &dns.A{
		Hdr: dns.RR_Header{
			Name:   q.Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    300,
		},
		A: blockedIP,
	}
	m.Answer = append(m.Answer, rr)

	s.mu.Lock()
	s.blockedRequests = append(s.blockedRequests, models.BlockedRequest{
		Domain: q.Name,
		IP:     blockedIP.String(),
		Source: sourceAddr,
	})
	s.mu.Unlock()
}

func (s *Server) Start() error {
	dns.HandleFunc(".", s.handleDNS)
	server := &dns.Server{Addr: ":53", Net: "udp"}
	return server.ListenAndServe()
}

func (s *Server) GetBlockedRequests() []models.BlockedRequest {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.BlockedRequest{}, s.blockedRequests...)
}
