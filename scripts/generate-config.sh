#!/bin/bash
set -e

CONFIG_FILE="backend/config/server.json"
OUTPUT_FORMAT="json"
GENERATE_QR=false

# Parse arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --qr) GENERATE_QR=true ;;
        --format) OUTPUT_FORMAT="$2"; shift ;;
        *) echo "Unknown parameter: $1"; exit 1 ;;
    esac
    shift
done

if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ Server config not found at $CONFIG_FILE"
    exit 1
fi

# Extract server details
SERVER_IP=$(jq -r '.locations[0].ip' "$CONFIG_FILE")
SERVER_PORT=$(jq -r '.locations[0].port' "$CONFIG_FILE")
PROTOCOL=$(jq -r '.locations[0].protocol' "$CONFIG_FILE")
ENCRYPTION=$(jq -r '.encryption' "$CONFIG_FILE")

# Generate auth token (simple example - use proper JWT in production)
AUTH_TOKEN=$(openssl rand -hex 32)

# Read CA certificate
CA_CERT=""
if [ -f "certs/ca.crt" ]; then
    CA_CERT=$(cat certs/ca.crt | sed ':a;N;$!ba;s/\n/\\n/g')
fi

CREATED_AT=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
EXPIRES_AT=$(date -u -d "+365 days" +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -v+365d +"%Y-%m-%dT%H:%M:%SZ")

# Generate client config JSON
CLIENT_CONFIG=$(cat <<EOF
{
  "server": "$SERVER_IP",
  "port": $SERVER_PORT,
  "protocol": "$PROTOCOL",
  "encryption": "$ENCRYPTION",
  "auth_token": "$AUTH_TOKEN",
  "ca_cert": "$CA_CERT",
  "created_at": "$CREATED_AT",
  "expires_at": "$EXPIRES_AT"
}
EOF
)

if [ "$GENERATE_QR" = true ]; then
    if ! command -v qrencode &> /dev/null; then
        echo "📦 Installing qrencode..."
        sudo apt-get install -y qrencode 2>/dev/null || brew install qrencode 2>/dev/null
    fi
    echo "$CLIENT_CONFIG" | qrencode -t ANSIUTF8
    echo ""
    echo "📱 Scan this QR code with the WarnetechVPN app"
else
    echo "$CLIENT_CONFIG"
fi
