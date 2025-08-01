package dns

import (
	"context"
	"time"

	"github.com/karmaKiller3352/Xray-core/common/errors"
	"github.com/karmaKiller3352/Xray-core/common/log"
	"github.com/karmaKiller3352/Xray-core/common/net"
	"github.com/karmaKiller3352/Xray-core/features/dns"
	"github.com/karmaKiller3352/Xray-core/features/dns/localdns"
)

// LocalNameServer is an wrapper over local DNS feature.
type LocalNameServer struct {
	client *localdns.Client
}

// QueryIP implements Server.
func (s *LocalNameServer) QueryIP(ctx context.Context, domain string, option dns.IPOption) (ips []net.IP, ttl uint32, err error) {
	start := time.Now()
	ips, ttl, err = s.client.LookupIP(domain, option)

	if len(ips) > 0 {
		errors.LogInfo(ctx, "Localhost got answer: ", domain, " -> ", ips)
		log.Record(&log.DNSLog{Server: s.Name(), Domain: domain, Result: ips, Status: log.DNSQueried, Elapsed: time.Since(start), Error: err})
	}

	return
}

// Name implements Server.
func (s *LocalNameServer) Name() string {
	return "localhost"
}

// NewLocalNameServer creates localdns server object for directly lookup in system DNS.
func NewLocalNameServer() *LocalNameServer {
	errors.LogInfo(context.Background(), "DNS: created localhost client")
	return &LocalNameServer{
		client: localdns.New(),
	}
}

// NewLocalDNSClient creates localdns client object for directly lookup in system DNS.
func NewLocalDNSClient(ipOption dns.IPOption) *Client {
	return &Client{server: NewLocalNameServer(), ipOption: &ipOption}
}
