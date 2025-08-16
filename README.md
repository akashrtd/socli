# SOCLI - Decentralized Terminal Microblogging

**SOCLI** is a privacy-first, decentralized, terminal-based microblogging platform built in Go. It enables ephemeral peer-to-peer communication using the libp2p networking stack. All content is stored only in memory and vanishes when the application exits, ensuring maximum user privacy.

<p align="center">
  <img src="docs/screenshots/socli_demo.gif" alt="SOCLI Demo"/>
</p>

## Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Running SOCLI](#running-socli)
- [Usage](#usage)
  - [User Interface](#user-interface)
  - [Commands](#commands)
  - [Keybindings](#keybindings)
- [Architecture](#architecture)
  - [Core Components](#core-components)
  - [P2P Networking](#p2p-networking)
- [Configuration](#configuration)
- [Security & Privacy](#security--privacy)
- [Development](#development)
  - [Building from Source](#building-from-source)
  - [Running Tests](#running-tests)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Decentralized & P2P:** No central servers. Communicate directly with peers using libp2p.
- **Ephemeral Posts:** All content exists only in memory. Nothing is persisted to disk (except your private key).
- **Terminal UI (TUI):** Rich, interactive terminal interface built with Charm's BubbleTea.
- **Real-time Messaging:** Instantly broadcast and receive messages via GossipSub.
- **Markdown Support:** Format your posts with Markdown, rendered beautifully in the terminal.
- **Dynamic Hashtag Subscriptions:** Join and leave topics (#hashtags) on the fly to see relevant posts.
- **Peer Discovery:** Automatically discover other SOCLI users on your local network (mDNS) and globally (DHT).
- **End-to-End Encryption:** Message payloads are encrypted using NaCl Box before being sent over the network.
- **Message Signing:** All posts are cryptographically signed for authenticity.
- **Privacy First:** Designed from the ground up to minimize data retention and maximize user anonymity.

## Quick Start

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.21 or higher)

### Installation

**Option 1: Download a Binary (Recommended)**

Pre-built binaries for various platforms are available on the [Releases](https://github.com/yourusername/socli/releases) page. Download the appropriate binary for your operating system and architecture.

**Option 2: Install with `go install`**

```bash
go install github.com/yourusername/socli@latest
```
*(Note: Replace `yourusername` with your actual GitHub username if hosting there)*

**Option 3: Build from Source**

```bash
git clone https://github.com/yourusername/socli.git
cd socli
go build -o socli .
```

### Running SOCLI

Simply run the executable:

```bash
./socli # or just `socli` if installed via `go install`
```

When you first run SOCLI, it will:
1. Generate a new cryptographic identity (`socli.key`) and store it locally.
2. Start listening for connections on a random TCP port.
3. Begin discovering peers using mDNS (local network) and DHT (global).
4. Launch the Terminal User Interface (TUI).

You should see your unique Peer ID and the initial "general" feed.

## Usage

### User Interface

SOCLI features a dual-pane TUI:
- **Main Feed (Left):** Displays posts from subscribed topics.
- **Information Panel (Right):** Shows connected peers and subscribed topics.
- **Status Bar (Bottom):** Displays your Peer ID and key controls.

### Commands

While in the compose view (press `c`), you can enter special commands prefixed with `/`.

- **`/subscribe <hashtag>`**: Joins a new topic to start receiving posts tagged with `#hashtag`.
- **`/unsubscribe <hashtag>`**: Leaves a topic to stop receiving posts for `#hashtag`.
- *(More commands will be added in future releases)*

### Keybindings

- **Global:**
  - `q` or `Ctrl+C`: Quit the application.
- **Feed View (Default):**
  - `c`: Switch to the compose view to write a new post or enter a command.
  - `j`: Scroll down to older posts.
  - `k`: Scroll up to newer posts.
- **Compose View:**
  - `Enter`: Send the typed message or execute the command.
  - `Esc`: Discard the current message/command and return to the feed view.

## Architecture

SOCLI is structured into several distinct Go packages, each responsible for a specific domain.

### Core Components

- **`main.go`**: Application entry point. Orchestrates the initialization of all core components (networking, storage, UI, messaging) and starts the main event loop.
- **`tui/`**: Terminal User Interface built with [BubbleTea](https://github.com/charmbracelet/bubbletea). Manages the application's state and view rendering.
- **`p2p/`**: Handles all libp2p networking, including peer discovery (mDNS/DHT), host creation, and pubsub setup.
- **`messaging/`**: Core messaging logic, including `Message` structure, topic management (`socli/hashtag/{tag}`), and the `Broadcaster`.
- **`storage/`**: Ephemeral in-memory data management for posts and discovered peers.
- **`crypto/`**: Cryptographic operations for signing messages and encrypting payloads using NaCl Box.
- **`content/`**: Content processing, primarily integrating [Glamour](https://github.com/charmbracelet/glamour) for Markdown rendering.
- **`config/`**: Loading and saving application configuration to `config.yaml`.
- **`internal/`**: Utility functions (e.g., hashtag extraction).

### P2P Networking

SOCLI leverages the robust [libp2p](https://libp2p.io/) stack for all networking functionalities:

1.  **Identity & Transport Security:**
    *   Each node generates a unique Ed25519 key pair (managed by libp2p) for its Peer ID.
    *   Connections between peers are secured using the Noise protocol.

2.  **Peer Discovery:**
    *   **Local Network (mDNS):** Automatically discovers other SOCLI nodes on the same local network segment.
    *   **Global Network (DHT):** Enables discovery of nodes across the internet using the Kademlia Distributed Hash Table.

3.  **Messaging (PubSub):**
    *   Uses `GossipSub` for efficient, scalable, and resilient real-time message broadcasting.
    *   Messages are routed based on topics. SOCLI uses the naming convention `socli/hashtag/{hashtag}` for topics.
    *   When you subscribe to a hashtag (e.g., `#tech`), SOCLI joins the `socli/hashtag/tech` topic.
    *   Publishing a post with hashtags causes it to be broadcast to *all* relevant topics simultaneously.

## Configuration

On first run, SOCLI creates a `config.yaml` file in its directory with default settings.

```yaml
network:
  listen_port: 0 # 0 means a random port, or specify a port (e.g., 4001)
  bootstrap_peers: [] # List of multiaddrs for bootstrap nodes (future feature)
  enable_mdns: true # Enable mDNS for local peer discovery
  enable_dht: true # Enable Kademlia DHT for global peer discovery
ui:
  theme: "default" # TUI theme (currently unused)
  refresh_rate_ms: 100 # UI refresh rate in milliseconds (currently unused)
  max_post_length: 280 # Maximum character length for posts
privacy:
  encrypt_messages: true # Enable end-to-end encryption for message payloads
  key_path: "socli.key" # Path to store the private key file
  auto_clear_on_exit: true # Automatically clear all in-memory data on exit
```

You can modify this file to customize network settings, UI preferences, and privacy features.

## Security & Privacy

SOCLI prioritizes user privacy and data security:

- **Ephemerality:** All posts and peer information are stored in volatile memory and are erased when the application closes.
- **Cryptographic Identity:** Each user/node has a unique libp2p Peer ID derived from a secret key, ensuring identity without a central authority.
- **Transport Security:** All direct connections between peers are encrypted using libp2p's Noise protocol.
- **Application-Layer Encryption:** Message payloads are further encrypted using NaCl Box before being published via pubsub. While the current implementation encrypts with the sender's own key (for simplicity), the framework allows for true E2E encryption for specific recipients in the future.
- **Message Integrity:** Every post is signed with the sender's private key, allowing recipients to verify authenticity.
- **Local Key Storage:** Your private key is stored locally in `socli.key` (configurable) and is never transmitted. *Protect this file.*
- **No Central Servers:** There are no third parties that can collect or analyze your data.

## Development

### Building from Source

1. Ensure you have Go installed (version specified in `go.mod`).
2. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/socli.git
   cd socli
   ```
3. Build the binary:
   ```bash
   go build -o socli .
   ```

### Running Tests

SOCLI uses Go's built-in testing framework. To run all tests:

```bash
go test ./...
```

To run tests for a specific package with verbose output:

```bash
cd p2p && go test -v
```

## Contributing

We welcome and appreciate contributions from the community! Whether it's bug reports, feature requests, code contributions, or documentation improvements, your help is invaluable.

See [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to contribute.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details.