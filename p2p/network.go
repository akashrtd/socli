package p2p

import (
	"context"
	"socli/config"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
)

// NetworkManager handles the libp2p host and networking functionality.
type NetworkManager struct {
	Host host.Host
	cfg  *config.Config
}

// NewNetworkManager creates and initializes a new libp2p host.
func NewNetworkManager(cfg *config.Config) (*NetworkManager, error) {
	// Create a new libp2p host
	h, err := libp2p.New(
		// Use the options constructor to configure the host
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Security(noise.ID, noise.New),
		libp2p.DefaultMuxers,
	)
	if err != nil {
		return nil, err
	}

	return &NetworkManager{
		Host: h,
		cfg:  cfg,
	}, nil
}

// Start begins the networking operations like peer discovery.
func (nm *NetworkManager) Start(ctx context.Context) error {
	if nm.cfg.Network.EnableMDNS {
		if err := setupMDNSDiscovery(ctx, nm.Host); err != nil {
			return err
		}
	}

	if nm.cfg.Network.EnableDHT {
		if _, err := setupDHTDiscovery(ctx, nm.Host); err != nil {
			return err
		}
	}

	return nil
}

// Close shuts down the libp2p host.
func (nm *NetworkManager) Close() error {
	return nm.Host.Close()
}
