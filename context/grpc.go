package context

import (
	"fmt"
	"sync"

	dg "github.com/kovey/discovery/grpc"
	"github.com/kovey/discovery/krpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	grpcDialOpts   []grpc.DialOption
	grpcDialOptsMu sync.RWMutex
)

// SetGRPCDialOptions sets the gRPC dial options used for all service connections.
// Use this to enable mTLS — for example:
//
//	creds, _ := credentials.NewClientTLSFromFile("ca.pem", "server.example.com")
//	context.SetGRPCDialOptions(grpc.WithTransportCredentials(creds))
//
// To enable mutual TLS, use a tls.Config with client certificates and wrap it
// with credentials.NewTLS.
//
// Pass nil or call with no arguments to restore insecure defaults.
func SetGRPCDialOptions(opts ...grpc.DialOption) {
	grpcDialOptsMu.Lock()
	defer grpcDialOptsMu.Unlock()
	if len(opts) == 0 {
		grpcDialOpts = nil
		return
	}
	grpcDialOpts = opts
}

func dialOptions() []grpc.DialOption {
	grpcDialOptsMu.RLock()
	defer grpcDialOptsMu.RUnlock()
	if len(grpcDialOpts) == 0 {
		return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	}
	opts := make([]grpc.DialOption, len(grpcDialOpts))
	copy(opts, grpcDialOpts)
	return opts
}

// dialGRPC creates a gRPC client connection for the given service.
// If custom dial options are set via SetGRPCDialOptions, they are used;
// otherwise insecure credentials are used.
func dialGRPC(serviceName krpc.ServiceName, group string) (grpc.ClientConnInterface, error) {
	target := fmt.Sprintf("%s://%s", dg.Scheme_Etcd, serviceName.Group(group))
	opts := dialOptions()
	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("grpc dial %s: %w", target, err)
	}
	return conn, nil
}
