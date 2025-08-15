package p2p

import (
	"context"
	"log"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

// discoveryNotifee handles peer discovery notifications.
// It now also sends a notification to the TUI.
type discoveryNotifee struct {
	h             host.Host
	peerConnected func(peer.ID) // Callback function to notify about connected peers
}

// HandlePeerFound connects to peers discovered via mDNS and notifies the application.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	log.Printf("Discovery: Found new peer: %s\n", pi.ID.String()) // Changed log message
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		log.Printf("Discovery: Error connecting to peer %s: %s\n", pi.ID.String(), err)
		return // Don't notify if connection failed
	}
	log.Printf("Discovery: Successfully connected to peer: %s\n", pi.ID.String()) // Add success log

	// Notify the application (e.g., TUI) about the newly connected peer
	if n.peerConnected != nil {
		n.peerConnected(pi.ID)
	}
}

// setupMDNSDiscovery initializes mDNS for local peer discovery.
// It accepts a callback for peer connection notifications.
func setupMDNSDiscovery(ctx context.Context, h host.Host, onPeerConnected func(peer.ID)) error {
	// setup mDNS discovery
	service := mdns.NewMdnsService(h, "socli-discovery", &discoveryNotifee{h: h, peerConnected: onPeerConnected})
	return service.Start()
}

// setupDHTDiscovery initializes the Kademlia DHT for global peer discovery.
// It accepts a callback for peer connection notifications.
func setupDHTDiscovery(ctx context.Context, h host.Host, onPeerConnected func(peer.ID)) (*dht.IpfsDHT, error) {
	// Start a DHT, for use in peer discovery.
	kademliaDHT, err := dht.New(ctx, h)
	if err != nil {
		return nil, err
	}

	// Bootstrap the DHT. In the default configuration, this spawns a background
	// thread that will refresh the peer table every five minutes.
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		return nil, err
	}

	// We're not connecting to any bootstrap peers here. In a real-world application,
	// you would want to connect to a list of bootstrap peers.

	return kademliaDHT, nil
}
