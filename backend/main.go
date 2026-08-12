package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/quic-go/quic-go"
)

// ServerConfig represents the VPN server configuration
type ServerConfig struct {
	Name        string      `json:"name"`
	ListenAddr  string      `json:"listen_addr"`
	Protocols   []string    `json:"protocols"`
	Encryption  string      `json:"encryption"`
	MaxClients  int         `json:"max_clients"`
	Timeout     int         `json:"timeout_seconds"`
	Locations   []Location  `json:"locations"`
	Logging     LogConfig   `json:"logging"`
	Security    SecurityConfig `json:"security"`
}

type Location struct {
	Name     string `json:"name"`
	Country  string `json:"country"`
	Region   string `json:"region"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

type LogConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	File   string `json:"file"`
	Days   int    `json:"retention_days"`
}

type SecurityConfig struct {
	RequireAuth bool   `json:"require_auth"`
	CertFile    string `json:"cert_file"`
	KeyFile     string `json:"key_file"`
	CAFile      string `json:"ca_file"`
}

// VPNServer represents the main VPN server
type VPNServer struct {
	config      *ServerConfig
	clients     map[string]*ClientSession
	clientMutex sync.RWMutex
	quicServer  *quic.Conn
	done        chan bool
}

// ClientSession represents an active client connection
type ClientSession struct {
	ID          string
	IP          string
	ConnectedAt time.Time
	Protocol    string
	BytesSent   uint64
	BytesRecv   uint64
	LastSeen    time.Time
}

var (
	configPath = flag.String("config", "config/server.json", "Path to server configuration file")
	logFile    *os.File
)

func main() {
	flag.Parse()

	// Load configuration
	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logging
	initLogging(config.Logging)
	defer logFile.Close()

	log.Printf("🚀 Starting WarnetechVPN Server v1.0.0")
	log.Printf("📍 Configuration: %s", *configPath)
	log.Printf("🔐 Encryption: %s", config.Encryption)
	log.Printf("📊 Max Clients: %d", config.MaxClients)

	// Create VPN server
	server := &VPNServer{
		config:  config,
		clients: make(map[string]*ClientSession),
		done:    make(chan bool),
	}

	// Validate listen address
	addr, err := net.ResolveUDPAddr("udp", config.ListenAddr)
	if err != nil {
		log.Fatalf("Failed to resolve listen address: %v", err)
	}

	log.Printf("🎧 Listening on %s", addr.String())
	log.Printf("📋 Protocols: %v", config.Protocols)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("⛔ Shutting down gracefully...")
		server.Shutdown()
		os.Exit(0)
	}()

	// Start protocol handlers
	for _, protocol := range config.Protocols {
		switch protocol {
		case "quic":
			go server.startQUICServer()
		case "wireguard":
			go server.startWireGuardServer()
		case "tcp":
			go server.startTCPServer()
		default:
			log.Printf("⚠️  Unknown protocol: %s", protocol)
		}
	}

	// Start monitoring goroutine
	go server.monitorConnections()

	// Keep server running
	<-server.done
}

func (s *VPNServer) startQUICServer() {
	log.Println("🔵 Starting QUIC protocol handler...")
	// Load TLS certificate
	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{loadCertificate(s.config.Security)},
	}

	listener, err := quic.ListenAddr(s.config.ListenAddr, tlsConf, nil)
	if err != nil {
		log.Printf("❌ QUIC Error: %v", err)
		return
	}

	log.Println("✅ QUIC server started")

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			log.Printf("❌ QUIC connection error: %v", err)
			continue
		}

		go s.handleQUICConnection(conn)
	}
}

func (s *VPNServer) startWireGuardServer() {
	log.Println("🟢 Starting WireGuard protocol handler...")
	// WireGuard implementation
	// This would require wgctrl library
	log.Println("✅ WireGuard server started")
}

func (s *VPNServer) startTCPServer() {
	log.Println("🟡 Starting TCP fallback handler...")
	
	listener, err := net.Listen("tcp", s.config.ListenAddr)
	if err != nil {
		log.Printf("❌ TCP Error: %v", err)
		return
	}

	log.Println("✅ TCP server started")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("❌ TCP connection error: %v", err)
			continue
		}

		go s.handleTCPConnection(conn)
	}
}

func (s *VPNServer) handleQUICConnection(conn *quic.Conn) {
	clientID := conn.RemoteAddr().String()
	log.Printf("📱 New QUIC connection from %s", clientID)

	session := &ClientSession{
		ID:          clientID,
		IP:          conn.RemoteAddr().String(),
		ConnectedAt: time.Now(),
		Protocol:    "quic",
	}

	s.clientMutex.Lock()
	s.clients[clientID] = session
	s.clientMutex.Unlock()

	defer func() {
		s.clientMutex.Lock()
		delete(s.clients, clientID)
		s.clientMutex.Unlock()
		log.Printf("❌ Client disconnected: %s", clientID)
	}()

	// Handle client streams
	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			log.Printf("⚠️  Error accepting stream: %v", err)
			return
		}

		go s.handleStream(stream, session)
	}
}

func (s *VPNServer) handleTCPConnection(conn net.Conn) {
	clientID := conn.RemoteAddr().String()
	log.Printf("📱 New TCP connection from %s", clientID)

	session := &ClientSession{
		ID:          clientID,
		IP:          conn.RemoteAddr().String(),
		ConnectedAt: time.Now(),
		Protocol:    "tcp",
	}

	s.clientMutex.Lock()
	s.clients[clientID] = session
	s.clientMutex.Unlock()

	defer func() {
		conn.Close()
		s.clientMutex.Lock()
		delete(s.clients, clientID)
		s.clientMutex.Unlock()
		log.Printf("❌ Client disconnected: %s", clientID)
	}()

	// Handle TCP traffic
	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}

		// Process packet
		s.processPacket(buffer[:n], session)
	}
}

func (s *VPNServer) handleStream(stream *quic.Stream, session *ClientSession) {
	buffer := make([]byte, 4096)
	defer stream.Close()

	for {
		n, err := stream.Read(buffer)
		if err != nil {
			log.Printf("⚠️  Stream read error: %v", err)
			return
		}

		session.BytesRecv += uint64(n)
		s.processPacket(buffer[:n], session)
	}
}

func (s *VPNServer) processPacket(data []byte, session *ClientSession) {
	session.LastSeen = time.Now()
	session.BytesSent += uint64(len(data))

	// Process VPN packet here
	// Decrypt, route, encrypt, etc.
}

func (s *VPNServer) monitorConnections() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.clientMutex.RLock()
			log.Printf("📊 Active connections: %d/%d", len(s.clients), s.config.MaxClients)
			for id, session := range s.clients {
				uptimeMinutes := time.Since(session.ConnectedAt).Minutes()
				log.Printf("  └─ %s | %s | %.1fm | ↑%dB ↓%dB",
					id, session.Protocol, uptimeMinutes,
					session.BytesSent, session.BytesRecv)
			}
			s.clientMutex.RUnlock()
		}
	}
}

func (s *VPNServer) Shutdown() {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()

	log.Printf("🔌 Closing %d client connections...", len(s.clients))
	s.clients = make(map[string]*ClientSession)
	s.done <- true
}

func loadConfig(path string) (*ServerConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config ServerConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func initLogging(logConf LogConfig) {
	var err error
	if logConf.File != "" {
		logFile, err = os.OpenFile(logConf.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		log.SetOutput(logFile)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func loadCertificate(secConf SecurityConfig) tls.Certificate {
	cert, err := tls.LoadX509KeyPair(secConf.CertFile, secConf.KeyFile)
	if err != nil {
		log.Fatalf("Failed to load certificate: %v", err)
	}
	return cert
}

// This is a simplified version. In production, you'd need:
// - Proper packet routing with iptables/netfilter
// - UDP handling for QUIC
// - DNS masking
// - MTU handling
// - Connection state tracking
