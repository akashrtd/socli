# AI Coding Agent Prompt: Build SOCLI - Decentralized P2P Social Platform

## Project Overview
Build **socli**, a privacy-first, decentralized, terminal-based social platform in Go that enables ephemeral peer-to-peer communication with Markdown support, ASCII art sharing, and real-time messaging through a rich TUI interface.

## Core Requirements

### Functional Requirements
1. **Terminal UI (TUI)**: Rich terminal interface using BubbleTea framework
2. **P2P Networking**: Decentralized peer-to-peer communication using libp2p
3. **Ephemeral Posts**: All content vanishes when application exits (no persistent storage)
4. **Markdown Support**: Full Markdown rendering in terminal using Glamour
5. **Real-time Messaging**: Live post broadcasting and receiving via GossipSub
6. **Hashtag System**: Topic-based post filtering and peer discovery
7. **Encryption**: Local private key encryption using NaCl Box
8. **ASCII Art**: Support for ASCII art in posts
9. **Privacy-First**: No central servers, no data persistence, local key storage only

### Technical Stack
- **Language**: Go (latest version)
- **TUI Framework**: github.com/charmbracelet/bubbletea + bubbles
- **P2P Network**: github.com/libp2p/go-libp2p + go-libp2p-pubsub
- **Markdown**: github.com/charmbracelet/glamour
- **Encryption**: golang.org/x/crypto/nacl/box
- **Discovery**: mDNS + DHT for peer discovery

## Architecture Requirements

### Core Components
1. **TUI Layer** (`tui/`): BubbleTea application with multiple views
   - Feed view (main posts timeline)
   - Compose view (create new posts)
   - Profile view (user settings)
   - Notifications view (real-time updates)

2. **P2P Network Layer** (`p2p/`): libp2p networking
   - Host initialization with security
   - GossipSub for topic-based messaging
   - mDNS and DHT for peer discovery
   - Custom protocol handlers

3. **Messaging System** (`messaging/`): Message handling
   - Post structure with metadata
   - Hashtag extraction and routing
   - Message broadcasting and filtering

4. **Encryption Layer** (`crypto/`): Security implementation
   - NaCl Box public-key encryption
   - Local key generation and storage
   - Message signing and verification

5. **Content Processing** (`content/`): Rendering and parsing
   - Glamour Markdown rendering
   - ASCII art processing
   - Content sanitization

6. **Ephemeral Storage** (`storage/`): In-memory data management
   - Post storage and retrieval
   - Peer information caching
   - Session state management

## Implementation Specifications

### Project Structure
```
socli/
├── main.go                 # Application entry point
├── go.mod                  # Go module with dependencies
├── go.sum                  # Dependency checksums
├── README.md               # Usage instructions
├── config/
│   ├── config.go          # Configuration management
│   └── defaults.go        # Default settings
├── tui/                   # Terminal UI components
│   ├── app.go             # Main BubbleTea application
│   ├── keys.go            # Keyboard shortcuts
│   ├── views/
│   │   ├── feed.go        # Main feed view
│   │   ├── compose.go     # Post composition
│   │   ├── profile.go     # User profile/settings
│   │   └── notifications.go # Real-time notifications
│   └── components/
│       ├── post.go        # Individual post display
│       ├── input.go       # Enhanced text input
│       └── sidebar.go     # Navigation sidebar
├── p2p/                   # P2P networking
│   ├── network.go         # Network manager
│   ├── discovery.go       # Peer discovery (mDNS/DHT)
│   ├── pubsub.go         # GossipSub message handling
│   ├── peer.go           # Peer management
│   └── protocols.go      # Custom protocol definitions
├── messaging/             # Message handling
│   ├── message.go        # Message structures
│   ├── topics.go         # Topic management
│   ├── broadcast.go      # Message broadcasting
│   └── filters.go        # Content filtering
├── crypto/                # Encryption layer
│   ├── keys.go           # Key management
│   ├── encryption.go     # Encryption/decryption
│   └── identity.go       # User identity
├── content/               # Content processing
│   ├── renderer.go       # Markdown rendering
│   ├── ascii.go          # ASCII art handling
│   └── sanitizer.go      # Content sanitization
├── storage/               # Ephemeral storage
│   ├── memory.go         # In-memory storage
│   ├── posts.go          # Post management
│   └── peers.go          # Peer information
└── internal/              # Internal utilities
    ├── utils.go          # Utility functions
    └── constants.go      # Application constants
```

### Key Dependencies (go.mod)
```go
module socli

go 1.21

require (
    github.com/charmbracelet/bubbletea v0.24.2
    github.com/charmbracelet/bubbles v0.16.1
    github.com/charmbracelet/glamour v0.6.0
    github.com/charmbracelet/lipgloss v0.8.0
    github.com/libp2p/go-libp2p v0.30.0
    github.com/libp2p/go-libp2p-pubsub v0.9.3
    github.com/libp2p/go-libp2p-kad-dht v0.24.2
    golang.org/x/crypto v0.12.0
    gopkg.in/yaml.v3 v3.0.1
)
```

### Core Data Structures

#### Message Structure
```go
type Message struct {
    ID        string    `json:"id"`
    Author    string    `json:"author"`    // Peer ID
    Content   string    `json:"content"`   // Markdown content
    Hashtags  []string  `json:"hashtags"`  // Extracted hashtags
    Timestamp time.Time `json:"timestamp"`
    Signature []byte    `json:"signature"` // Message signature
    Type      MsgType   `json:"type"`      // Post, Reply, Share
    ReplyTo   string    `json:"reply_to,omitempty"`
}
```

#### Network Configuration
```go
type Config struct {
    Network struct {
        ListenPort     int      `yaml:"listen_port"`
        BootstrapPeers []string `yaml:"bootstrap_peers"`
        EnableMDNS     bool     `yaml:"enable_mdns"`
        EnableDHT      bool     `yaml:"enable_dht"`
    } `yaml:"network"`
    
    UI struct {
        Theme         string `yaml:"theme"`
        RefreshRate   int    `yaml:"refresh_rate_ms"`
        MaxPostLength int    `yaml:"max_post_length"`
    } `yaml:"ui"`
    
    Privacy struct {
        EncryptMessages bool   `yaml:"encrypt_messages"`
        KeyPath        string `yaml:"key_path"`
        AutoClear      bool   `yaml:"auto_clear_on_exit"`
    } `yaml:"privacy"`
}
```

## Implementation Requirements

### 1. BubbleTea Application Structure
- Implement main App model with Update/View methods
- Create multiple view states (feed, compose, profile, notifications)
- Handle keyboard navigation and shortcuts
- Responsive layout that adapts to terminal size
- Real-time updates without blocking UI

### 2. libp2p Network Implementation
- Initialize host with proper security and multiplexing
- Set up GossipSub for topic-based messaging
- Implement mDNS for local peer discovery
- Configure DHT for global peer discovery
- Handle connection management and error recovery

### 3. Message Broadcasting System
- Subscribe to hashtag-based topics (`socli/hashtag/{hashtag}`)
- Broadcast messages to relevant topics
- Handle message encryption/decryption
- Implement message validation and filtering
- Support direct peer-to-peer messaging

### 4. Encryption and Security
- Generate Ed25519 key pairs for identity
- Use NaCl Box for message encryption
- Implement message signing for authenticity
- Store keys encrypted locally with passphrase
- No plaintext storage of private data

### 5. Content Rendering
- Integrate Glamour for Markdown rendering
- Support for code blocks, tables, lists, emphasis
- ASCII art preservation and display
- Hashtag parsing and highlighting
- Content sanitization for security

### 6. User Interface Features
- Vim-style navigation (j/k for scrolling)
- Tab completion for hashtags and commands
- Live Markdown preview in compose view
- Real-time peer count and status
- Notification indicators for new posts

### 7. Ephemeral Data Management
- All posts stored only in memory
- Automatic cleanup on application exit
- No persistent storage of user content
- Session-based peer information caching
- Optional export functionality for important posts

## Specific Implementation Details

### libp2p Host Initialization
```go
host, err := libp2p.New(
    libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
    libp2p.Security(noise.ID, noise.New),
    libp2p.Muxer(yamux.ID, yamux.DefaultTransport),
    libp2p.DefaultTransports,
)
```

### GossipSub Topic Management
```go
// Subscribe to hashtag topics
topicName := fmt.Sprintf("socli/hashtag/%s", hashtag)
topic, err := ps.Join(topicName)
subscription, err := topic.Subscribe()
```

### BubbleTea Message Types
```go
type PostReceivedMsg struct{ Post Message }
type PeerConnectedMsg struct{ PeerID string }
type NetworkErrorMsg struct{ Err error }
type WindowResizeMsg struct{ Width, Height int }
```

### Encryption Implementation
```go
// Encrypt message for recipient
encrypted := box.Seal(nonce[:], messageBytes, &nonce, recipientPubKey, senderPrivKey)

// Decrypt received message  
decrypted, ok := box.Open(nil, encrypted[24:], &nonce, senderPubKey, recipientPrivKey)
```

## Expected Deliverables

### 1. Working Application
- Complete Go application that compiles and runs
- All core features functional (posting, receiving, rendering)
- Error handling and graceful degradation
- Proper resource cleanup on exit

### 2. Documentation
- README.md with installation and usage instructions
- Code comments explaining complex logic
- Configuration file documentation
- Architecture overview comments

### 3. Testing
- Unit tests for core functionality
- Integration tests for P2P communication
- Example posts demonstrating Markdown features
- Performance benchmarks for message handling

### 4. Configuration
- Default configuration file (config.yaml)
- Command-line argument parsing
- Environment variable support
- Configurable key bindings

## Quality Requirements

### Performance
- Handle 1000+ concurrent peers efficiently
- Message rendering under 100ms
- Startup time under 5 seconds
- Memory usage under 100MB for normal operation

### Security
- No plaintext private key storage
- Message integrity verification
- Protection against spam/flooding
- Secure peer authentication

### Usability
- Intuitive keyboard shortcuts
- Clear visual feedback for actions
- Error messages with helpful guidance
- Responsive interface design

### Code Quality
- Follow Go best practices and conventions
- Proper error handling throughout
- Clean, readable code structure
- Comprehensive logging for debugging

## Success Criteria

The implementation is successful when:

1. **Users can run the application** and connect to other peers automatically
2. **Posts are created, encrypted, and broadcast** to the P2P network in real-time
3. **Markdown content renders beautifully** in the terminal interface
4. **Hashtag-based discovery works** for finding relevant posts and peers
5. **Data is truly ephemeral** - no traces left after application exit
6. **Interface is responsive** and enjoyable to use in terminal
7. **Network is resilient** to peer disconnections and network issues

Build this as a complete, production-ready application that demonstrates the power of decentralized, privacy-first social networking in the terminal environment.