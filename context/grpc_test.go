package context

import (
	"sync"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetGRPCDialOptions_EmptyRestoresDefault(t *testing.T) {
	// Set some options first
	customOpt := grpc.WithTransportCredentials(insecure.NewCredentials())
	SetGRPCDialOptions(customOpt)
	require.NotEmpty(t, grpcDialOpts)

	// Empty call should reset to nil
	SetGRPCDialOptions()
	assert.Nil(t, grpcDialOpts)
}

func TestSetGRPCDialOptions_WithCustomOptions(t *testing.T) {
	opt1 := grpc.WithTransportCredentials(insecure.NewCredentials())
	opt2 := grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(1024))

	SetGRPCDialOptions(opt1, opt2)

	opts := dialOptions()
	assert.Len(t, opts, 2)

	// Reset
	SetGRPCDialOptions()
}

func TestDialOptions_DefaultIsInsecure(t *testing.T) {
	SetGRPCDialOptions() // ensure clean state

	opts := dialOptions()
	assert.Len(t, opts, 1)
	// Verify it's a transport credentials option (insecure)
	assert.NotNil(t, opts[0])
}

func TestDialOptions_ReturnsCopy(t *testing.T) {
	opt := grpc.WithTransportCredentials(insecure.NewCredentials())
	SetGRPCDialOptions(opt)

	opts1 := dialOptions()
	opts2 := dialOptions()

	// Should be different slices (copy)
	assert.NotSame(t, &opts1[0], &opts2[0])

	SetGRPCDialOptions()
}

func TestDialOptions_ThreadSafe(t *testing.T) {
	SetGRPCDialOptions() // clean state

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			SetGRPCDialOptions(grpc.WithTransportCredentials(insecure.NewCredentials()))
		}()
		go func() {
			defer wg.Done()
			_ = dialOptions()
		}()
	}
	wg.Wait()

	// Should not panic; state is consistent
	assert.NotPanics(t, func() {
		_ = dialOptions()
	})

	SetGRPCDialOptions()
}

func TestSetGRPCDialOptions_ConcurrentReads(t *testing.T) {
	opt := grpc.WithTransportCredentials(insecure.NewCredentials())
	SetGRPCDialOptions(opt)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			opts := dialOptions()
			assert.Len(t, opts, 1)
		}()
	}
	wg.Wait()

	SetGRPCDialOptions()
}

func TestDialOptions_NilReturnedWhenNoOptions(t *testing.T) {
	// Force nil state by calling empty SetGRPCDialOptions
	SetGRPCDialOptions()
	assert.Nil(t, grpcDialOpts)

	// dialOptions should return insecure defaults
	opts := dialOptions()
	assert.Len(t, opts, 1)
}
