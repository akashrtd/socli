# SOCLI Encryption Strategy

This document outlines the current encryption strategy employed by SOCLI and discusses potential future enhancements.

## Current Implementation

SOCLI implements a hybrid approach to security, leveraging both libp2p's built-in transport security and application-layer payload encryption.

### 1. Transport Layer Security (TLS)

All direct peer-to-peer connections established by libp2p are automatically secured using the **Noise Protocol**. This ensures that data in transit between any two nodes is encrypted and authenticated, protecting against eavesdropping and man-in-the-middle attacks at the network level.

### 2. Application-Layer Payload Encryption

Beyond transport security, SOCLI adds an additional layer by encrypting the **payload** of each `messaging.Message` before it is published via GossipSub. This is handled in `messaging/broadcaster.go`.

- **Mechanism:** Uses `golang.org/x/crypto/nacl/box` (NaCl Box) for symmetric encryption.
- **Key Pair:** A dedicated Curve25519 key pair (`crypto.KeyPair`) is generated and managed by SOCLI, separate from the libp2p host's identity key. This key pair is used for both encryption and decryption of payloads.
- **Process:**
    1. When a user sends a post, the `messaging.Broadcaster` serializes the `Message` struct (JSON).
    2. If `config.Privacy.EncryptMessages` is `true`, the serialized JSON data is encrypted using `crypto.Encrypt`.
    3. `crypto.Encrypt` uses the **sender's own `PrivateKey` and `PublicKey`** to perform the encryption.
    4. The resulting encrypted blob is then published to all relevant GossipSub topics.

### 3. Payload Decryption

When a message is received:

- The `main.go` subscription loop (or dynamic subscription handlers in `tui/subscription.go`) receives the raw `[]byte` data from GossipSub.
- If `config.Privacy.EncryptMessages` is `true`, it attempts to decrypt the data using `crypto.Decrypt`.
- `crypto.Decrypt` uses the **receiver's own `PrivateKey` and `PublicKey`** (the same key pair used for encryption).
- If decryption is successful, the original JSON is recovered and unmarshalled into a `messaging.Message`.
- If decryption fails (e.g., the data wasn't encrypted by the receiver's key, or wasn't encrypted at all), the process logs a warning and attempts to parse the raw data as unencrypted JSON (to maintain backward compatibility or handle messages from nodes with encryption off).

## Rationale for Current Approach

The choice to encrypt the payload with the sender's own key was a pragmatic one for the prototype:

- **Simplicity:** It provides an additional obfuscation layer without the complexity of managing recipient public keys in a broadcast pubsub model.
- **Consistency:** Every node uses its own key pair for both encrypting outbound messages and attempting to decrypt inbound messages.
- **Baseline Privacy:** It ensures that message content is not plaintext on the wire, adding another hurdle beyond transport security.

## Limitations & Future Considerations

The current approach has a significant limitation: it is **not true end-to-end encryption for specific recipients**. Because the same key is used to encrypt and decrypt, any node using SOCLI with the same key (which is every node, as they all use their own key) would theoretically be able to decrypt the message. However, in practice, a node can only successfully decrypt messages that were encrypted with *its own* key.

A more robust E2E encryption scheme for specific peers would involve:

1.  **Asymmetric Key Exchange:** Nodes would need to exchange their public `crypto.KeyPair.PublicKey`.
2.  **Per-Recipient Encryption:** When broadcasting a message, the sender would encrypt the payload *multiple times*, once for each intended recipient's public key (or use a group encryption scheme).
3.  **Key Management:** A system for securely storing and retrieving the public keys of peers with whom you wish to communicate privately.

This is a non-trivial extension, especially in a fully decentralized discovery model, and is considered a future enhancement beyond the scope of the initial prototype.

## Conclusion

The current encryption strategy provides a solid baseline of privacy by ensuring message payloads are not plaintext and by leveraging strong transport security. It acknowledges the complexities of implementing true E2E encryption in a broadcast pubsub system and provides a functional stepping stone. Future work will focus on evolving this strategy to support true E2E communication for private messaging or private channels.