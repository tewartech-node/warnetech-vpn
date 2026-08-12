#!/bin/bash
set -e

echo "🚀 WarnetechVPN Server Setup"
echo "=============================="

# Check for Docker
if ! command -v docker &> /dev/null; then
    echo "📦 Docker not found. Installing..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    rm get-docker.sh
fi

if ! command -v docker-compose &> /dev/null; then
    echo "📦 Docker Compose not found. Installing..."
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
fi

# Generate certificates if they don't exist
if [ ! -f "certs/server.crt" ]; then
    echo "🔐 Generating TLS certificates..."
    mkdir -p certs
    ./scripts/generate-certs.sh
fi

# Get server's public IP
PUBLIC_IP=$(curl -s https://api.ipify.org)
echo "🌐 Detected public IP: $PUBLIC_IP"

# Update config with detected IP
if [ -f "backend/config/server.json" ]; then
    sed -i.bak "s/203.0.113.1/$PUBLIC_IP/g" backend/config/server.json
    echo "✅ Updated server config with your public IP"
fi

# Open firewall ports
if command -v ufw &> /dev/null; then
    echo "🔥 Configuring firewall..."
    sudo ufw allow 443/tcp
    sudo ufw allow 443/udp
    sudo ufw allow 51820/udp
    echo "✅ Firewall configured"
fi

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "  1. Review config: backend/config/server.json"
echo "  2. Start server: docker-compose up -d"
echo "  3. Generate client config: ./scripts/generate-config.sh"
echo "  4. Share with friends: ./scripts/share-config.sh --create-link"
echo ""
