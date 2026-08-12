#!/bin/bash
set -e

CERT_DIR="certs"
DAYS_VALID=825

mkdir -p "$CERT_DIR"
cd "$CERT_DIR"

echo "🔐 Generating CA certificate..."
openssl genrsa -out ca.key 4096
openssl req -x509 -new -nodes -key ca.key -sha256 -days $DAYS_VALID \
    -out ca.crt \
    -subj "/C=US/ST=State/L=City/O=WarnetechVPN/CN=WarnetechVPN-CA"

echo "🔐 Generating server certificate..."
openssl genrsa -out server.key 4096
openssl req -new -key server.key -out server.csr \
    -subj "/C=US/ST=State/L=City/O=WarnetechVPN/CN=vpn.warnetech.local"

openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key \
    -CAcreateserial -out server.crt -days $DAYS_VALID -sha256

rm server.csr

echo "✅ Certificates generated in $CERT_DIR/"
echo "  - ca.crt / ca.key      (Certificate Authority)"
echo "  - server.crt / server.key  (Server certificate)"
echo ""
echo "⚠️  Keep .key files secure and never commit them to git!"
