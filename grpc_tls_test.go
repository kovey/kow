package kow

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func generateTestCert(t *testing.T, dir, name string) (certFile, keyFile string) {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName: name,
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	require.NoError(t, err)

	certPath := filepath.Join(dir, name+".pem")
	keyPath := filepath.Join(dir, name+"-key.pem")

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	require.NoError(t, os.WriteFile(certPath, certPEM, 0644))

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	require.NoError(t, os.WriteFile(keyPath, keyPEM, 0600))

	return certPath, keyPath
}

func TestSetGRPCTLS_NilConfig(t *testing.T) {
	// First set a TLS config to ensure it gets cleared
	cfg := &GRPCTLSConfig{InsecureSkipVerify: true}
	require.NoError(t, SetGRPCTLS(cfg))
	assert.True(t, IsGRPCTLSEnabled())

	// Nil config should disable TLS and reset to insecure
	require.NoError(t, SetGRPCTLS(nil))
	assert.False(t, IsGRPCTLSEnabled())
	assert.Nil(t, grpcTLSConfig)
	assert.Nil(t, grpcTLSCreds)
}

func TestSetGRPCTLS_InsecureSkipVerify(t *testing.T) {
	cfg := &GRPCTLSConfig{
		ServerName:         "test.example.com",
		InsecureSkipVerify: true,
	}
	require.NoError(t, SetGRPCTLS(cfg))
	assert.True(t, IsGRPCTLSEnabled())
	assert.NotNil(t, grpcTLSCreds)

	require.NoError(t, SetGRPCTLS(nil))
}

func TestSetGRPCTLS_WithCACert(t *testing.T) {
	dir := t.TempDir()
	caCert, _ := generateTestCert(t, dir, "ca")

	cfg := &GRPCTLSConfig{
		CACertFile: caCert,
	}
	require.NoError(t, SetGRPCTLS(cfg))
	assert.True(t, IsGRPCTLSEnabled())

	require.NoError(t, SetGRPCTLS(nil))
}

func TestSetGRPCTLS_WithClientCert(t *testing.T) {
	dir := t.TempDir()
	certFile, keyFile := generateTestCert(t, dir, "client")

	cfg := &GRPCTLSConfig{
		CertFile: certFile,
		KeyFile:  keyFile,
	}
	require.NoError(t, SetGRPCTLS(cfg))
	assert.True(t, IsGRPCTLSEnabled())

	require.NoError(t, SetGRPCTLS(nil))
}

func TestSetGRPCTLS_FullMTLS(t *testing.T) {
	dir := t.TempDir()
	caCert, _ := generateTestCert(t, dir, "ca")
	clientCert, clientKey := generateTestCert(t, dir, "client")

	cfg := &GRPCTLSConfig{
		ServerName: "grpc.example.com",
		CACertFile: caCert,
		CertFile:   clientCert,
		KeyFile:    clientKey,
	}
	require.NoError(t, SetGRPCTLS(cfg))
	assert.True(t, IsGRPCTLSEnabled())
	assert.NotNil(t, grpcTLSCreds)
	assert.Equal(t, cfg, grpcTLSConfig)

	require.NoError(t, SetGRPCTLS(nil))
}

func TestSetGRPCTLS_InvalidCACert(t *testing.T) {
	dir := t.TempDir()
	badPath := filepath.Join(dir, "nonexistent.pem")

	cfg := &GRPCTLSConfig{
		CACertFile: badPath,
	}
	err := SetGRPCTLS(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "read ca cert")
}

func TestSetGRPCTLS_InvalidClientCert(t *testing.T) {
	dir := t.TempDir()
	badPath := filepath.Join(dir, "nonexistent.pem")
	keyFile := filepath.Join(dir, "nonexistent-key.pem")

	cfg := &GRPCTLSConfig{
		CertFile: badPath,
		KeyFile:  keyFile,
	}
	err := SetGRPCTLS(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "load client cert")
}

func TestSetGRPCTLS_InvalidCAPEM(t *testing.T) {
	dir := t.TempDir()
	badCert := filepath.Join(dir, "bad.pem")
	require.NoError(t, os.WriteFile(badCert, []byte("not a valid pem"), 0644))

	cfg := &GRPCTLSConfig{
		CACertFile: badCert,
	}
	err := SetGRPCTLS(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parse ca cert")
}

func TestSetGRPCTLS_MissingKeyFile(t *testing.T) {
	dir := t.TempDir()
	certFile, _ := generateTestCert(t, dir, "client")

	// CertFile without KeyFile — no client cert loaded, but no error
	cfg := &GRPCTLSConfig{
		CertFile: certFile,
	}
	require.NoError(t, SetGRPCTLS(cfg))
	assert.True(t, IsGRPCTLSEnabled())

	require.NoError(t, SetGRPCTLS(nil))
}

func TestSetGRPCTLSFromEnv(t *testing.T) {
	dir := t.TempDir()
	caCert, _ := generateTestCert(t, dir, "ca")
	clientCert, clientKey := generateTestCert(t, dir, "client")

	t.Setenv("GRPC_TLS_SERVER_NAME", "grpc.example.com")
	t.Setenv("GRPC_TLS_CA_CERT_FILE", caCert)
	t.Setenv("GRPC_TLS_CLIENT_CERT_FILE", clientCert)
	t.Setenv("GRPC_TLS_CLIENT_KEY_FILE", clientKey)

	require.NoError(t, SetGRPCTLSFromEnv())
	assert.True(t, IsGRPCTLSEnabled())

	require.NoError(t, SetGRPCTLS(nil))
}

func TestSetGRPCTLSFromEnv_NoVars(t *testing.T) {
	require.NoError(t, SetGRPCTLSFromEnv())
	assert.False(t, IsGRPCTLSEnabled())
}

func TestSetGRPCTLSFromEnv_PartialVars(t *testing.T) {
	dir := t.TempDir()
	caCert, _ := generateTestCert(t, dir, "ca")

	t.Setenv("GRPC_TLS_CA_CERT_FILE", caCert)
	// Other vars not set

	require.NoError(t, SetGRPCTLSFromEnv())
	assert.True(t, IsGRPCTLSEnabled())

	require.NoError(t, SetGRPCTLS(nil))
}

func TestIsGRPCTLSEnabled_Default(t *testing.T) {
	require.NoError(t, SetGRPCTLS(nil))
	assert.False(t, IsGRPCTLSEnabled())
}

func TestBuildGRPCCredentials_ServerTLS(t *testing.T) {
	dir := t.TempDir()
	caCert, _ := generateTestCert(t, dir, "ca")

	cfg := &GRPCTLSConfig{
		ServerName: "myserver",
		CACertFile: caCert,
	}

	creds, err := buildGRPCCredentials(cfg)
	require.NoError(t, err)
	assert.NotNil(t, creds)
	assert.NotNil(t, creds.Info())
}

func TestBuildGRPCCredentials_MutualTLS(t *testing.T) {
	dir := t.TempDir()
	caCert, _ := generateTestCert(t, dir, "ca")
	clientCert, clientKey := generateTestCert(t, dir, "client")

	cfg := &GRPCTLSConfig{
		ServerName: "mtls-server",
		CACertFile: caCert,
		CertFile:   clientCert,
		KeyFile:    clientKey,
	}

	creds, err := buildGRPCCredentials(cfg)
	require.NoError(t, err)
	assert.NotNil(t, creds)
}

func TestSetGRPCTLS_Overwrite(t *testing.T) {
	dir := t.TempDir()
	ca1, _ := generateTestCert(t, dir, "ca1")
	ca2, _ := generateTestCert(t, dir, "ca2")

	require.NoError(t, SetGRPCTLS(&GRPCTLSConfig{CACertFile: ca1}))
	assert.True(t, IsGRPCTLSEnabled())

	require.NoError(t, SetGRPCTLS(&GRPCTLSConfig{CACertFile: ca2}))
	assert.True(t, IsGRPCTLSEnabled())

	require.NoError(t, SetGRPCTLS(nil))
	assert.False(t, IsGRPCTLSEnabled())
}

func TestSetGRPCTLS_CloneCredentials(t *testing.T) {
	dir := t.TempDir()
	caCert, _ := generateTestCert(t, dir, "ca")

	cfg := &GRPCTLSConfig{
		ServerName: "clone-test",
		CACertFile: caCert,
	}

	require.NoError(t, SetGRPCTLS(cfg))
	creds := grpcTLSCreds
	require.NotNil(t, creds)

	cloned := creds.Clone()
	assert.NotNil(t, cloned)

	require.NoError(t, SetGRPCTLS(nil))
}
