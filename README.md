# SOCLI - Decentralized P2P Social Platform

**SOCLI** is a privacy-first, decentralized, terminal-based social platform built in Go. It enables ephemeral peer-to-peer communication with rich Markdown support, ASCII art sharing, and real-time messaging through an intuitive Terminal User Interface (TUI). Our mission is to provide a secure, private, and censorship-resistant social experience without relying on centralized servers or persistent data storage.

## Project Overview

SOCLI leverages modern P2P technologies to create a robust and privacy-focused social platform. All content is ephemeral, vanishing when the application exits, ensuring maximum user privacy. Communication is topic-based, allowing users to discover and interact with like-minded peers through a hashtag system.

## Core Requirements

### Functional Requirements

1.  **Terminal UI (TUI)**: Rich and interactive terminal interface using the Charm BubbleTea framework.
2.  **P2P Networking**: Decentralized peer-to-peer communication powered by `libp2p`.
3.  **Ephemeral Posts**: All content (posts, messages) is stored only in memory and vanishes upon application exit, ensuring no persistent user data.
4.  **Markdown Support**: Full Markdown rendering within the terminal using `Glamour` for rich text formatting.
5.  **Real-time Messaging**: Live post broadcasting and receiving via `GossipSub` for efficient topic-based communication.
6.  **Hashtag System**: Topic-based post filtering and peer discovery, allowing users to subscribe to specific interests.
7.  **Encryption**: Local private key encryption using `NaCl Box` for secure message exchange.
8.  **ASCII Art**: Support for displaying ASCII art within posts.
9.  **Privacy-First**: No central servers, no data persistence, and local key storage only.

### Technical Stack

*   **Language**: Go (latest stable version, currently Go 1.21+)
*   **TUI Framework**: `github.com/charmbracelet/bubbletea` + `github.com/charmbracelet/bubbles`
*   **P2P Network**: `github.com/libp2p/go-libp2p` + `github.com/libp2p/go-libp2p-pubsub`
*   **Markdown Rendering**: `github.com/charmbracelet/glamour`
*   **Terminal Styling**: `github.com/charmbracelet/lipgloss`
*   **Encryption**: `golang.org/x/crypto/nacl/box`
*   **Discovery**: mDNS (local) + Kademlia DHT (global) for peer discovery.
*   **Configuration**: `gopkg.in/yaml.v3` for YAML-based configuration.

## Architecture

SOCLI is structured into several distinct packages, each responsible for a specific domain:

```
socli/
├── main.go                 # Application entry point, initializes and wires up core components.
├── go.mod                  # Go module dependencies.
├── go.sum                  # Dependency checksums.
├── README.md               # Project documentation (this file).
├── config/                 # Handles application configuration loading and default settings.
│   ├── config.go           # Defines the Config struct and LoadConfig function.
│   └── defaults.go         # Provides default configuration values.
├── tui/                    # Terminal User Interface components using BubbleTea.
│   ├── app.go              # Main BubbleTea application model, manages views and state.
│   ├── keys.go             # Defines keyboard shortcuts and key bindings.
│   ├── views/              # Different screens/views of the application.
│   │   ├── feed.go         # Displays the main posts timeline.
│   │   ├── compose.go      # Allows users to create and send new posts.
│   │   ├── profile.go      # User profile and settings view (placeholder).
│   │   └── notifications.go# Real-time notifications view (placeholder).
│   └── components/         # Reusable UI elements.
│       ├── post.go         # Component for displaying individual posts (placeholder).
│       ├── input.go        # Enhanced text input component (e.g., for compose view).
│       └── sidebar.go      # Navigation sidebar component (placeholder).
├── p2p/                    # libp2p networking layer.
│   ├── network.go          # Manages the libp2p host, connections, and lifecycle.
│   ├── discovery.go        # Handles mDNS and DHT for peer discovery.
│   ├── pubsub.go           # Manages GossipSub for topic-based messaging.
│   ├── peer.go             # Peer management utilities (placeholder).
│   └── protocols.go        # Custom libp2p protocol definitions (placeholder).
├── messaging/              # Handles message structures, broadcasting, and filtering.
│   ├── message.go          # Defines the core Message struct.
│   ├── topics.go           # Manages topic naming conventions (e.g., for hashtags).
│   ├── broadcast.go        # Logic for broadcasting messages to the network.
│   └── filters.go          # Content filtering and validation (placeholder).
├── crypto/                 # Encryption and cryptographic operations.
│   ├── keys.go             # Handles key pair generation, loading, and saving.
│   ├── encryption.go       # Implements message encryption/decryption using NaCl Box.
│   └── identity.go         # User identity management (placeholder).
├── content/                # Content processing, rendering, and sanitization.
│   ├── renderer.go         # Integrates Glamour for Markdown rendering.
│   ├── ascii.go            # Handles ASCII art processing (placeholder).
│   └── sanitizer.go        # Content sanitization for security (placeholder).
├── storage/                # Ephemeral in-memory data management.
│   ├── memory.go           # Implements the in-memory store for posts and peers.
│   ├── posts.go            # Post management within the store (placeholder).
│   └── peers.go            # Peer information management within the store (placeholder).
└── internal/               # Internal utilities and application constants.
    ├── utils.go            # General utility functions (placeholder).
    └── constants.go        # Defines application-wide constants.
```

## Getting Started

To get started with SOCLI, you'll need to have Go (version 1.21 or higher) installed on your system.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/akashrtd/socli.git
    cd socli
    ```
2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```
    This command will download all the necessary Go modules and update `go.sum`.

3.  **Build the application (optional):**
    ```bash
    go build -o socli .
    ```
    This will create an executable named `socli` in your current directory.

4.  **Run the application:**
    ```bash
    go run .
    ```
    Alternatively, if you built the executable:
    ```bash
    ./socli
    ```
    This will launch the SOCLI application in your terminal. You should see your Peer ID and listening addresses.

## Usage

Once the application is running:

*   **Quit**: Press `q` or `Ctrl+C` to exit the application.
*   **Compose**: Press `c` to switch to the compose view. Type your message and press `Enter` to send. Press `Esc` to return to the feed view without sending.
*   **Feed**: The main view displays incoming messages from other peers.

## Configuration

SOCLI uses a `config.yaml` file for its settings. A default `config.yaml` will be created in the application's root directory if it doesn't exist.

Example `config.yaml`:

```yaml
network:
  listen_port: 0          # 0 means a random port, or specify a port (e.g., 4001)
  bootstrap_peers: []     # List of multiaddrs for bootstrap nodes (e.g., "/ip4/1.2.3.4/tcp/4001/p2p/Qm...")
  enable_mdns: true       # Enable mDNS for local peer discovery
  enable_dht: true        # Enable Kademlia DHT for global peer discovery
ui:
  theme: "default"        # TUI theme (e.g., "dark", "light")
  refresh_rate_ms: 100    # UI refresh rate in milliseconds
  max_post_length: 280    # Maximum character length for posts
privacy:
  encrypt_messages: true  # Enable end-to-end encryption for messages
  key_path: "socli.key"   # Path to store the private key file
  auto_clear_on_exit: true # Automatically clear all in-memory data on exit
```

You can modify this file to customize network settings, UI preferences, and privacy features.

## Key Bindings

*   `q` / `Ctrl+C`: Quit the application.
*   `c`: Switch to compose view.
*   `Enter` (in compose view): Send message.
*   `Esc` (in compose view): Return to feed view without sending.

## Contributing

We welcome and appreciate contributions from the community! Whether it's bug reports, feature requests, code contributions, or documentation improvements, your help is invaluable.

### How to Contribute

1.  **Fork the Repository**: Start by forking the `socli` repository on GitHub.
2.  **Clone Your Fork**:
    ```bash
    git clone https://github.com/YOUR_USERNAME/socli.git
    cd socli
    ```
3.  **Create a New Branch**:
    ```bash
    git checkout -b feature/your-feature-name
    # or
    git checkout -b bugfix/issue-description
    ```
    Choose a descriptive branch name.
4.  **Install Dependencies**: Ensure all Go modules are up-to-date.
    ```bash
    go mod tidy
    ```
5.  **Make Your Changes**: Implement your feature or bug fix.
    *   Adhere to Go best practices and coding conventions.
    *   Ensure your code is clean, readable, and well-commented.
    *   Run `gofmt` and `golint` (or `golangci-lint`) to ensure code style consistency.
6.  **Write Tests**:
    *   All new features and bug fixes should be accompanied by appropriate unit and/or integration tests.
    *   Tests should cover edge cases and ensure the stability of the application.
7.  **Run Tests**:
    ```bash
    go test ./...
    ```
    Ensure all tests pass before submitting your changes.
8.  **Commit Your Changes**: Write clear, concise commit messages.
    ```bash
    git commit -m "feat: Add new feature for X"
    # or
    git commit -m "fix: Resolve issue with Y"
    ```
9.  **Push to Your Fork**:
    ```bash
    git push origin feature/your-feature-name
    ```
10. **Submit a Pull Request (PR)**:
    *   Go to the original `socli` repository on GitHub.
    *   You should see a prompt to create a new pull request from your branch.
    *   Provide a detailed description of your changes, why they were made, and any relevant issue numbers.
    *   Ensure your PR passes all CI checks.

### Code Style and Quality

*   **Go Best Practices**: Follow idiomatic Go programming patterns.
*   **Readability**: Write code that is easy to understand and maintain.
*   **Error Handling**: Implement robust error handling throughout the codebase.
*   **Modularity**: Keep functions and packages focused on single responsibilities.
*   **Documentation**: Add comments for complex logic and public functions/types.

### Reporting Issues

If you find a bug or have a feature request, please open an issue on our [GitHub Issues page](https://github.com/akashrtd/socli/issues). Provide as much detail as possible, including steps to reproduce bugs, expected behavior, and your environment details.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details.

## Contact

For any questions, feedback, or discussions, please open an issue on our [GitHub Issues page](https://github.com/akashrtd/socli/issues).
