package p2p

import (
	"context"
	"socli/config"
	"testing"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/connmgr"
	"github.com/libp2p/go-libp2p/core/event"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"
)

// mockHost is a simple mock implementation of host.Host for testing.
type mockHost struct {
	id peer.ID
}

func (m *mockHost) ID() peer.ID {
	return m.id
}

func (m *mockHost) Peerstore() peerstore.Peerstore {
	// For testing purposes, we can return nil or a mock peerstore.
	// The NetworkManager.Start method doesn't use Peerstore directly.
	return nil
}

func (m *mockHost) Addrs() []ma.Multiaddr {
	// For testing purposes, we can return an empty slice.
	// The NetworkManager.Start method doesn't use Addrs directly.
	return []ma.Multiaddr{}
}

func (m *mockHost) Network() network.Network {
	// For testing purposes, we can return nil.
	// The NetworkManager.Start method doesn't use Network directly.
	return nil
}

func (m *mockHost) Mux() protocol.Switch {
	// For testing purposes, we can return nil.
	// The NetworkManager.Start method doesn't use Mux directly.
	return nil
}

func (m *mockHost) Connect(ctx context.Context, pi peer.AddrInfo) error {
	// For testing purposes, we can return nil to simulate success.
	// The NetworkManager.Start method doesn't call Connect directly.
	return nil
}

func (m *mockHost) SetStreamHandler(pid protocol.ID, handler network.StreamHandler) {
	// For testing purposes, we can do nothing.
	// The NetworkManager.Start method doesn't call SetStreamHandler.
}

func (m *mockHost) SetStreamHandlerMatch(pid protocol.ID, matcher func(protocol.ID) bool, handler network.StreamHandler) {
	// For testing purposes, we can do nothing.
	// The NetworkManager.Start method doesn't call SetStreamHandlerMatch.
}

func (m *mockHost) RemoveStreamHandler(pid protocol.ID) {
	// For testing purposes, we can do nothing.
	// The NetworkManager.Start method doesn't call RemoveStreamHandler.
}

func (m *mockHost) NewStream(ctx context.Context, p peer.ID, pids ...protocol.ID) (network.Stream, error) {
	// For testing purposes, we can return nil, nil.
	// The NetworkManager.Start method doesn't call NewStream.
	return nil, nil
}

func (m *mockHost) Close() error {
	// For testing purposes, we can return nil.
	// The NetworkManager.Start method doesn't call Close.
	return nil
}

func (m *mockHost) ConnManager() connmgr.ConnManager {
	// For testing purposes, we can return nil.
	// The NetworkManager.Start method doesn't use ConnManager.
	return nil
}

func (m *mockHost) EventBus() event.Bus {
	// For testing purposes, we can return nil.
	// The NetworkManager.Start method doesn't use EventBus.
	return nil
}

// mockSetupMDNSDiscovery is a mock implementation of setupMDNSDiscoveryFunc for testing.
func mockSetupMDNSDiscovery(ctx context.Context, h host.Host, onPeerConnected func(peer.ID)) error {
	// In a real test, we might record that this function was called.
	// For now, we'll just return nil to simulate success.
	return nil
}

// mockSetupDHTDiscovery is a mock implementation of setupDHTDiscoveryFunc for testing.
func mockSetupDHTDiscovery(ctx context.Context, h host.Host, onPeerConnected func(peer.ID)) (*dht.IpfsDHT, error) {
	// In a real test, we might record that this function was called and return a mock DHT.
	// For now, we'll just return nil, nil to simulate success.
	return nil, nil
}

// TestNetworkManagerStart tests the Start method of NetworkManager.
func TestNetworkManagerStart(t *testing.T) {
	// Create a mock libp2p host
	mockHost := &mockHost{id: peer.ID("test-host-id")}

	// Create a dummy config
	cfg := config.DefaultConfig()

	// Variables to track if the mock functions are called
	mdnsCalled := false
	dhtCalled := false

	// Create mock setup functions that record if they are called
	mockSetupMDNS := func(ctx context.Context, h host.Host, onPeerConnected func(peer.ID)) error {
		mdnsCalled = true
		return nil
	}

	mockSetupDHT := func(ctx context.Context, h host.Host, onPeerConnected func(peer.ID)) (*dht.IpfsDHT, error) {
		dhtCalled = true
		return nil, nil
	}

	// Create a NetworkManager with the mock host and mock setup functions
	nm := &NetworkManager{
		Host:      mockHost,
		cfg:       cfg,
		setupMDNS: mockSetupMDNS,
		setupDHT:  mockSetupDHT,
	}

	// Test with both MDNS and DHT enabled (default config)
	ctx := context.Background()
	err := nm.Start(ctx)
	if err != nil {
		t.Errorf("Start() with MDNS and DHT enabled returned error: %v, want nil", err)
	}

	// Verify that both setup functions were called
	if !mdnsCalled {
		t.Error("setupMDNSDiscovery was not called when MDNS was enabled")
	}
	if !dhtCalled {
		t.Error("setupDHTDiscovery was not called when DHT was enabled")
	}

	// Reset flags for next test
	mdnsCalled = false
	dhtCalled = false

	// Test with MDNS disabled and DHT enabled
	cfg.Network.EnableMDNS = false
	cfg.Network.EnableDHT = true
	nm.cfg = cfg
	err = nm.Start(ctx)
	if err != nil {
		t.Errorf("Start() with MDNS disabled and DHT enabled returned error: %v, want nil", err)
	}

	// Verify that only DHT setup function was called
	if mdnsCalled {
		t.Error("setupMDNSDiscovery was called when MDNS was disabled")
	}
	if !dhtCalled {
		t.Error("setupDHTDiscovery was not called when DHT was enabled")
	}

	// Reset flags for next test
	mdnsCalled = false
	dhtCalled = false

	// Test with MDNS enabled and DHT disabled
	cfg.Network.EnableMDNS = true
	cfg.Network.EnableDHT = false
	nm.cfg = cfg
	err = nm.Start(ctx)
	if err != nil {
		t.Errorf("Start() with MDNS enabled and DHT disabled returned error: %v, want nil", err)
	}

	// Verify that only MDNS setup function was called
	if !mdnsCalled {
		t.Error("setupMDNSDiscovery was not called when MDNS was enabled")
	}
	if dhtCalled {
		t.Error("setupDHTDiscovery was called when DHT was disabled")
	}

	// Reset flags for next test
	mdnsCalled = false
	dhtCalled = false

	// Test with both MDNS and DHT disabled
	cfg.Network.EnableMDNS = false
	cfg.Network.EnableDHT = false
	nm.cfg = cfg
	err = nm.Start(ctx)
	if err != nil {
		t.Errorf("Start() with both MDNS and DHT disabled returned error: %v, want nil", err)
	}

	// Verify that neither setup function was called
	if mdnsCalled {
		t.Error("setupMDNSDiscovery was called when both MDNS and DHT were disabled")
	}
	if dhtCalled {
		t.Error("setupDHTDiscovery was called when both MDNS and DHT were disabled")
	}
}