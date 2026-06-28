package kow

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/kovey/kow/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// GRPCTLSConfig holds the mTLS configuration for gRPC client connections.
type GRPCTLSConfig struct {
	// ServerName overrides the server name used for TLS verification (SNI).
	ServerName string

	// CACertFile is the path to the CA certificate file for server verification.
	// When empty, the host's root CA set is used (InsecureSkipVerify must be false).
	CACertFile string

	// CertFile is the path to the client certificate file for mutual TLS.
	CertFile string

	// KeyFile is the path to the client private key file for mutual TLS.
	KeyFile string

	// InsecureSkipVerify disables server certificate verification.
	// Only use for testing.
	InsecureSkipVerify bool
}

var (
	grpcTLSConfig *GRPCTLSConfig
	grpcTLSCreds  credentials.TransportCredentials
)

// SetGRPCTLS configures mTLS for all gRPC client connections.
// Call this before starting the server or at init time.
func SetGRPCTLS(cfg *GRPCTLSConfig) error {
	if cfg == nil {
		grpcTLSConfig = nil
		grpcTLSCreds = nil
		context.SetGRPCDialOptions()
		return nil
	}

	creds, err := buildGRPCCredentials(cfg)
	if err != nil {
		return fmt.Errorf("grpc tls: %w", err)
	}

	grpcTLSConfig = cfg
	grpcTLSCreds = creds

	// Apply TLS credentials to all future gRPC connections
	context.SetGRPCDialOptions(grpc.WithTransportCredentials(creds))
	return nil
}

// SetGRPCTLSFromEnv configures mTLS from environment variables:
//
//	GRPC_TLS_SERVER_NAME       - override server name (SNI)
//	GRPC_TLS_CA_CERT_FILE      - path to CA certificate file
//	GRPC_TLS_CLIENT_CERT_FILE  - path to client certificate file
//	GRPC_TLS_CLIENT_KEY_FILE   - path to client private key file
func SetGRPCTLSFromEnv() error {
	cfg := &GRPCTLSConfig{
		ServerName: os.Getenv("GRPC_TLS_SERVER_NAME"),
		CACertFile: os.Getenv("GRPC_TLS_CA_CERT_FILE"),
		CertFile:   os.Getenv("GRPC_TLS_CLIENT_CERT_FILE"),
		KeyFile:    os.Getenv("GRPC_TLS_CLIENT_KEY_FILE"),
	}
	if cfg.ServerName == "" && cfg.CACertFile == "" && cfg.CertFile == "" && cfg.KeyFile == "" {
		return nil
	}
	return SetGRPCTLS(cfg)
}

// IsGRPCTLSEnabled reports whether gRPC mTLS is configured.
func IsGRPCTLSEnabled() bool {
	return grpcTLSConfig != nil
}

func buildGRPCCredentials(cfg *GRPCTLSConfig) (credentials.TransportCredentials, error) {
	tlsCfg := &tls.Config{
		ServerName:         cfg.ServerName,
		InsecureSkipVerify: cfg.InsecureSkipVerify,
		MinVersion:         tls.VersionTLS12,
	}

	// Load CA certificate for server verification
	if cfg.CACertFile != "" {
		caCert, err := os.ReadFile(cfg.CACertFile)
		if err != nil {
			return nil, fmt.Errorf("read ca cert: %w", err)
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse ca cert from %s", cfg.CACertFile)
		}
		tlsCfg.RootCAs = caCertPool
	}

	// Load client certificate for mutual TLS
	if cfg.CertFile != "" && cfg.KeyFile != "" {
		clientCert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("load client cert: %w", err)
		}
		tlsCfg.Certificates = []tls.Certificate{clientCert}
	}

	return credentials.NewTLS(tlsCfg), nil
}
