# WarnetechVPN 🔐

A modern, user-friendly personal VPN application for Android with multi-protocol support (QUIC, WireGuard, TCP), designed for easy setup and instant sharing with friends.

## ✨ Features

- **Multi-Protocol Support**: QUIC (latest), WireGuard, TCP tunneling
- **Android Native App**: Beautiful Material Design 3 UI
- **Zero-Configuration**: Auto-detection and optimization
- **Easy Sharing**: QR codes, links, and direct imports
- **Lightweight**: Minimal battery drain, optimized data usage
- **Secure**: End-to-end encryption, ChaCha20 + AES-256
- **Kill Switch**: Automatic disconnect on VPN failure
- **Connection Monitor**: Real-time speed and statistics
- **Open Source**: Full transparency and auditable

## 🚀 Quick Start (Choose One)

### Fastest: Use Public Server
```bash
1. Download APK from Releases
2. Scan QR code from docs/quick-configs/
3. Tap Connect - Done!
```

### Self-Hosted Server
```bash
git clone https://github.com/tewartech-node/warnetech-vpn.git
cd warnetech-vpn
chmod +x scripts/*.sh
./scripts/setup-server.sh
docker-compose up -d
```

### Share with Friends
```bash
# Generate QR code config
./scripts/generate-config.sh --qr

# Create temporary 24h link
./scripts/share-config.sh --create-link

# They scan QR or click link, import config, and connect!
```

## 📋 Repository Structure

```
warnetech-vpn/
├── backend/
│   ├── main.go                 # Server entry point
│   ├── go.mod & go.sum         # Go dependencies
│   ├── server/
│   │   ├── quic.go            # QUIC protocol handler
│   │   ├── wireguard.go        # WireGuard support
│   │   └── tcp.go             # TCP fallback
│   ├── crypto/
│   │   └── encryption.go       # Crypto utilities
│   ├── config/
│   │   ├── server.json        # Server config template
│   │   └── config.go          # Config parser
│   └── Dockerfile
│
├── android/
│   ├── app/
│   │   ├── src/main/
│   │   │   ├── java/com/warnetech-vpn/
│   │   │   │   ├── MainActivity.kt
│   │   │   │   ├── VPNService.kt
│   │   │   │   ├── ConfigManager.kt
│   │   │   │   └── ConnectionMonitor.kt
│   │   │   ├── res/
│   │   │   │   ├── layout/
│   │   │   │   ├── values/
│   │   │   │   └── drawable/
│   │   │   └── AndroidManifest.xml
│   │   └── build.gradle.kts
│   ├── build.gradle.kts
│   ├── settings.gradle.kts
│   └── gradle.properties
│
├── scripts/
│   ├── setup-server.sh         # One-command server setup
│   ├── generate-config.sh      # Generate shareable configs
│   ├── share-config.sh         # Create sharing links
│   └── generate-certs.sh       # Certificate generation
│
├── docs/
│   ├── SETUP.md               # Detailed setup guide
│   ├── ANDROID-SETUP.md       # Android app build guide
│   ├── CONFIGURATION.md       # Configuration reference
│   ├── SHARING.md             # Share with friends guide
│   └── TROUBLESHOOTING.md     # Common issues
│
├── .github/workflows/
│   ├── backend-build.yml      # Go build & test
│   ├── android-build.yml      # APK build & sign
│   └── deploy.yml             # Auto-deploy on release
│
├── docker-compose.yml         # Local deployment
├── Dockerfile                 # Server container
├── .gitignore
├── LICENSE                    # MIT License
└── CONTRIBUTING.md            # Contribution guidelines
```

## 🔧 Installation & Setup

### Prerequisites

**Server Side:**
- Linux (Ubuntu 20.04+, Debian 11+) or any Docker-capable system
- 1GB RAM minimum (2GB+ recommended)
- 10GB storage
- Port 443 open (TCP/UDP)
- Optional: Port 51820 for WireGuard

**Android:**
- Android 10+ (API 29+)
- 100MB free storage
- VPN permission

**Development:**
- Go 1.19+
- Android Studio 2022.1+ (for building)
- Docker & Docker Compose

### Server Deployment

#### Option A: Docker (Easiest)
```bash
git clone https://github.com/tewartech-node/warnetech-vpn.git
cd warnetech-vpn

# Edit configuration (optional)
nano backend/config/server.json

# Start server
docker-compose up -d

# Verify it's running
docker-compose logs -f vpn-server

# Generate client config
./scripts/generate-config.sh > client-config.json
```

#### Option B: Manual Setup
```bash
git clone https://github.com/tewartech-node/warnetech-vpn.git
cd warnetech-vpn/backend

# Install Go dependencies
go mod download

# Generate certificates
../scripts/generate-certs.sh

# Build server
go build -o vpn-server main.go

# Run server
./vpn-server --config=../config/server.json
```

### Android App Installation

#### Option 1: Pre-built APK
1. Download latest APK from [Releases](https://github.com/tewartech-node/warnetech-vpn/releases)
2. Enable "Unknown Sources" in Android settings
3. Install APK
4. Open app → "Add Configuration" → Paste config JSON
5. Tap "Connect"

#### Option 2: Build from Source
```bash
cd android
./gradlew build
# APK located at: app/build/outputs/apk/release/app-release.apk
```

## 📱 How to Use the App

1. **Launch App**: Open WarnetechVPN
2. **Add Server Config**: 
   - Paste JSON config, OR
   - Scan QR code, OR
   - Import from link
3. **Select Server**: Choose from available locations
4. **Connect**: Tap the big connect button
5. **Monitor**: Watch real-time speed and data usage

## 🤝 Share with Friends (Easy Methods)

### Method 1: QR Code (Fastest)
```bash
./scripts/generate-config.sh --qr
# Shows QR code in terminal
# Friends scan with camera or VPN app → Auto-import!
```

### Method 2: Shareable Link (Safest)
```bash
./scripts/share-config.sh --create-link --ttl 24h
# Returns: https://vpn.example.com/share/abc123xyz
# Link expires after 24 hours
# Friends click → App auto-imports config
```

### Method 3: Export JSON (Manual)
```bash
./scripts/generate-config.sh
# Copy output → Share via Telegram, Discord, etc.
# Friends paste into app → Import
```

### Method 4: Direct App Import
```bash
# Email or message the generated config.json file
# Friends open VPN app → "Import Configuration" → Select file
```

## 🔐 Configuration Reference

### Server Config (`backend/config/server.json`)

```json
{
  "name": "WarnetechVPN",
  "listen_addr": "0.0.0.0:443",
  "protocols": ["quic", "wireguard", "tcp"],
  "encryption": "chacha20",
  "cipher_suites": ["AES-256-GCM", "ChaCha20-Poly1305"],
  "max_clients": 100,
  "timeout_seconds": 300,
  "locations": [
    {
      "name": "US East",
      "country": "US",
      "region": "us-east-1",
      "ip": "203.0.113.1",
      "port": 443,
      "protocol": "quic"
    },
    {
      "name": "EU Central",
      "country": "DE",
      "region": "eu-central-1",
      "ip": "203.0.113.2",
      "port": 443,
      "protocol": "wireguard"
    }
  ],
  "logging": {
    "level": "info",
    "format": "json",
    "file": "/var/log/warnetech-vpn.log",
    "retention_days": 7
  },
  "security": {
    "require_auth": true,
    "cert_file": "certs/server.crt",
    "key_file": "certs/server.key",
    "ca_file": "certs/ca.crt"
  }
}
```

### Client Config (Generated Automatically)

```json
{
  "server": "vpn.example.com",
  "port": 443,
  "protocol": "quic",
  "encryption": "chacha20",
  "auth_token": "eyJhbGciOiJIUzI1NiIs...",
  "ca_cert": "-----BEGIN CERTIFICATE-----\n...",
  "created_at": "2024-01-15T10:30:00Z",
  "expires_at": "2025-01-15T10:30:00Z"
}
```

## 📊 Supported Protocols

### QUIC (Default - Fastest)
- Latest protocol, optimized for mobile
- Faster handshakes
- Better for unstable connections
- Lower latency

### WireGuard (Secure)
- Lightweight and modern
- Excellent performance
- Simpler configuration
- Recommended for desktop/servers

### TCP Fallback
- Works everywhere (port 443)
- Slightly slower
- Automatic fallback if others fail

## 🛡️ Security & Privacy

- **Encryption**: ChaCha20-Poly1305 and AES-256-GCM
- **Authentication**: Certificate + Token-based
- **No Logging**: By default, no connection data stored
- **Kill Switch**: Automatic disconnect on VPN failure
- **Open Source**: Full code audit transparency
- **Regular Updates**: Security patches within 24 hours

### What We Don't Log
- ✅ No IP addresses
- ✅ No connection timestamps
- ✅ No bandwidth usage
- ✅ No DNS queries
- ✅ No application names

## 🚀 Advanced Features

### Connection Monitoring
```
Real-time stats:
- Upload/Download speed
- Latency (ping)
- Data transferred
- Active connections
- Server location
```

### Kill Switch
When enabled, if VPN disconnects unexpectedly, all internet access is blocked until you manually reconnect. Prevents data leaks.

### Split Tunneling
Route specific apps through VPN while others use direct connection (optional, for advanced users).

## 🔄 CI/CD & Automation

All commits trigger:
- ✅ Backend unit tests
- ✅ Android lint & build
- ✅ Security scanning
- ✅ Auto-release APK
- ✅ Docker image build

## 📚 Documentation

- [Setup Guide](docs/SETUP.md) - Detailed installation
- [Android Build Guide](docs/ANDROID-SETUP.md) - Build from source
- [Configuration Reference](docs/CONFIGURATION.md) - All config options
- [Sharing Guide](docs/SHARING.md) - Share with friends
- [Troubleshooting](docs/TROUBLESHOOTING.md) - Common issues

## 🤝 Contributing

Contributions are welcome!

1. Fork the repo
2. Create feature branch: `git checkout -b feature/amazing-feature`
3. Test thoroughly
4. Submit pull request

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📦 Releases

Latest: v1.0.0 (January 2024)
- APK downloads
- Docker image
- Source code releases

## 📄 License

MIT License - See [LICENSE](LICENSE)

Free to use, modify, and distribute. See license file for full terms.

## 🆘 Support & Community

- 📖 [Full Documentation](docs/)
- 🐛 [Report Issues](https://github.com/tewartech-node/warnetech-vpn/issues)
- 💬 [Discussions](https://github.com/tewartech-node/warnetech-vpn/discussions)
- 📧 Email: support@warnetech-vpn.com

## 🗺️ Roadmap

- [x] Multi-protocol support
- [x] Android app
- [ ] iOS support
- [ ] Windows/Mac desktop clients
- [ ] Speed test integration
- [ ] Analytics dashboard
- [ ] Built-in server marketplace
- [ ] Automated server setup wizard

## 🙏 Acknowledgments

Built with Go, Kotlin, and ❤️ for privacy advocates everywhere.

---

**Start using WarnetechVPN today. Total privacy, zero complexity.**
