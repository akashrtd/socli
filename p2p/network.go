package p2p

import (
	"context"
	"socli/config"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
)

// NetworkManager handles the libp2p host and networking functionality.
type NetworkManager struct {
	Host host.Host
	cfg  *config.Config
	// onPeerConnected is a callback function to notify when a new peer is connected.
	// This is set by the application (e.g., in main.go) to link discovery to the app logic (like TUI).
	onPeerConnected func(peer.ID) 
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
		// onPeerConnected will be set later by the application
	}, nil
}

// SetPeerConnectedCallback sets the callback function to be called when a peer is connected.
func (nm *NetworkManager) SetPeerConnectedCallback(callback func(peer.ID)) {
	nm.onPeerConnected = callback
}

// Start begins the networking operations like peer discovery.
func (nm *NetworkManager) Start(ctx context.Context) error {
	if nm.cfg.Network.EnableMDNS {
		if err := setupMDNSDiscovery(ctx, nm.Host, nm.onPeerConnected); err != nil {
			return err
		}
	}

	if nm.cfg.Network.EnableDHT {
		if _, err := setupDHTDiscovery(ctx, nm.Host, nm.onPeerConnected); err != nil {
			return err
		}
	}

	return nil
}

// Close shuts down the libp2p host.
func (nm *NetworkManager) Close() error {
	return nm.Host.Close()
}
