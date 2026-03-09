# Kendrick Lawton - homebrew-tap

The official Homebrew repository for **Kendrick Lawton (@k-henry)**. This tap serves as the centralized distribution point for **[Liquid Metal](https://liquidmetal.dev)** and other projects developed along the way.

## Installation

You can install tools directly from this tap using their fully qualified names:

```bash
brew install kendricklawton/tap/flux

```

For a better experience, add the tap to your local Homebrew installation to access all current and future formulae by their short names:

```bash
# Add the tap
brew tap kendricklawton/tap

# Install tools
brew install flux

```

## Available Formulae

This registry tracks the stable releases of the Liquid Metal ecosystem and standalone utilities developed for Arch Linux, Firecracker, and Wasm environments.

| Formula | Description |
| --- | --- |
| **`flux`** | **The Liquid Metal CLI** — The primary interface for shipping Firecracker microVMs and Wasm modules. |
| **`tba`** | *Upcoming projects, including custom eBPF monitors and specialized dev-tooling, will be indexed here.* |

---

## Architecture & Integration

All formulae in this tap are cross-compiled for **macOS (Intel/Apple Silicon)** and **Linux (amd64/arm64)**.

* **Automation**: Managed via GoReleaser.
* **Security**: Binaries are checksummed and verified against the `checksums.txt` provided in the GitHub Releases of the source repositories.
* **Updates**: To ensure you are running the latest version of any tool in this tap:
```bash
brew update && brew upgrade <formula>

```
