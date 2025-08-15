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
type discoveryNotifee struct {
	h host.Host
}

// HandlePeerFound connects to peers discovered via mDNS.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	log.Printf("Discovered a new peer: %s\n", pi.ID.String())
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		log.Printf("Error connecting to peer %s: %s\n", pi.ID.String(), err)
	}
}

// setupMDNSDiscovery initializes mDNS for local peer discovery.
func setupMDNSDiscovery(ctx context.Context, h host.Host) error {
	// setup mDNS discovery
	service := mdns.NewMdnsService(h, "socli-discovery", &discoveryNotifee{h: h})
	return service.Start()
}

// setupDHTDiscovery initializes the Kademlia DHT for global peer discovery.
func setupDHTDiscovery(ctx context.Context, h host.Host) (*dht.IpfsDHT, error) {
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
