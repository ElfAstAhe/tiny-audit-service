package grpc

import (
	"context"
	"crypto/x509"
	"testing"
	"time"

	mocks2 "github.com/ElfAstAhe/go-service-template/pkg/logger/mocks"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client"
	"github.com/ElfAstAhe/tiny-auth-service/pkg/transport/auth/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewOptions_Defaults(t *testing.T) {
	opts := NewOptions[string]()

	assert.NotNil(t, opts.BaseAuditClientOptions)
	assert.Equal(t, DefaultTarget, opts.Target)
	assert.Equal(t, DefaultConnectTimeout, opts.ConnTimeout)
	assert.Equal(t, DefaultIdleTimeout, opts.IdleTimeout)
	assert.Equal(t, DefaultKeepAliveTime, opts.KeepAliveTime)
	assert.Equal(t, DefaultKeepAliveTimeout, opts.KeepAliveTimeout)
	assert.Equal(t, DefaultKeepAlivePermitWithoutStream, opts.KeepAlivePermitWithoutStream)

	assert.Equal(t, DefaultWorkerCount, opts.Pool.WorkerCount)
	assert.Equal(t, DefaultDataCapacity, opts.Pool.DataCapacity)
	assert.Equal(t, DefaultCompleteProcessing, opts.Pool.CompleteProcess)
	assert.Equal(t, DefaultStopTimeout, opts.Pool.StopTimeout)
}

func TestOptions_FluentAPI(t *testing.T) {
	mockLogger := mocks2.NewMockLogger(t)
	mockProvider := mocks.NewMockTokenProvider(t)
	certPool := x509.NewCertPool()
	var dummyAction client.AuditAction[string] = func(ctx context.Context, workerIndex int, data string, token string) error { return nil }

	opts := NewOptions[string]()

	WithName[string]("grpc-audit-pool")(opts)
	WithTarget[string]("127.0.0.1:9090")(opts)
	WithCertPool[string](certPool)(opts)
	WithIdleTimeout[string](time.Minute)(opts)
	WithConnectTimeout[string](time.Second * 5)(opts)
	WithKeepAliveTime[string](time.Second * 30)(opts)
	WithKeepAliveTimeout[string](time.Second * 10)(opts)
	WithKeepAlivePermitWithoutStream[string](false)(opts)
	WithWorkerCount[string](8)(opts)
	WithDataCapacity[string](50000)(opts)
	WithCompleteProcess[string](false)(opts)
	WithStopTimeout[string](time.Second * 4)(opts)
	WithTokenProvider[string](mockProvider)(opts)
	withAuditAction[string](dummyAction)(opts)
	WithLogger[string](mockLogger)(opts)

	assert.Equal(t, "grpc-audit-pool", opts.Pool.Name)
	assert.Equal(t, "127.0.0.1:9090", opts.Target)
	assert.Equal(t, certPool, opts.CertPool)
	assert.Equal(t, time.Minute, opts.IdleTimeout)
	assert.Equal(t, time.Second*5, opts.ConnTimeout)
	assert.Equal(t, time.Second*30, opts.KeepAliveTime)
	assert.Equal(t, time.Second*10, opts.KeepAliveTimeout)
	assert.False(t, opts.KeepAlivePermitWithoutStream)
	assert.Equal(t, 8, opts.Pool.WorkerCount)
	assert.Equal(t, 50000, opts.Pool.DataCapacity)
	assert.False(t, opts.Pool.CompleteProcess)
	assert.Equal(t, time.Second*4, opts.Pool.StopTimeout)
	assert.Equal(t, mockProvider, opts.TokenProvider)
	assert.NotNil(t, opts.AuditAction)
	assert.Equal(t, mockLogger, opts.Pool.Logger)
}

func TestOptions_Validate(t *testing.T) {
	mockLogger := mocks2.NewMockLogger(t)
	mockProvider := mocks.NewMockTokenProvider(t)
	var dummyAction client.AuditAction[string] = func(ctx context.Context, workerIndex int, data string, token string) error { return nil }

	tests := []struct {
		name    string
		modify  func(*Options[string])
		wantErr bool
	}{
		{
			name: "valid options",
			modify: func(o *Options[string]) {
				o.Pool.Name = "valid-pool"
				o.Pool.Logger = mockLogger
				o.TokenProvider = mockProvider
				o.AuditAction = dummyAction
			},
			wantErr: false,
		},
		{
			name: "empty pool name",
			modify: func(o *Options[string]) {
				o.Pool.Name = "   "
				o.Pool.Logger = mockLogger
				o.TokenProvider = mockProvider
				o.AuditAction = dummyAction
			},
			wantErr: true,
		},
		{
			name: "invalid worker count",
			modify: func(o *Options[string]) {
				o.Pool.Name = "valid-pool"
				o.Pool.WorkerCount = 0
				o.Pool.Logger = mockLogger
				o.TokenProvider = mockProvider
				o.AuditAction = dummyAction
			},
			wantErr: true,
		},
		{
			name: "invalid data capacity",
			modify: func(o *Options[string]) {
				o.Pool.Name = "valid-pool"
				o.Pool.DataCapacity = -10
				o.Pool.Logger = mockLogger
				o.TokenProvider = mockProvider
				o.AuditAction = dummyAction
			},
			wantErr: true,
		},
		{
			name: "invalid stop timeout",
			modify: func(o *Options[string]) {
				o.Pool.Name = "valid-pool"
				o.Pool.StopTimeout = 0
				o.Pool.Logger = mockLogger
				o.TokenProvider = mockProvider
				o.AuditAction = dummyAction
			},
			wantErr: true,
		},
		{
			name: "nil pool logger",
			modify: func(o *Options[string]) {
				o.Pool.Name = "valid-pool"
				o.Pool.Logger = nil
				o.TokenProvider = mockProvider
				o.AuditAction = dummyAction
			},
			wantErr: true,
		},
		{
			name: "nil token provider",
			modify: func(o *Options[string]) {
				o.Pool.Name = "valid-pool"
				o.Pool.Logger = mockLogger
				o.TokenProvider = nil
				o.AuditAction = dummyAction
			},
			wantErr: true,
		},
		{
			name: "nil audit action",
			modify: func(o *Options[string]) {
				o.Pool.Name = "valid-pool"
				o.Pool.Logger = mockLogger
				o.TokenProvider = mockProvider
				o.AuditAction = nil
			},
			wantErr: true,
		},
		{
			name: "empty target",
			modify: func(o *Options[string]) {
				o.Pool.Name = "valid-pool"
				o.Target = ""
				o.Pool.Logger = mockLogger
				o.TokenProvider = mockProvider
				o.AuditAction = dummyAction
			},
			wantErr: true,
		},
		{
			name: "invalid connect timeout",
			modify: func(o *Options[string]) {
				o.Pool.Name = "valid-pool"
				o.ConnTimeout = 0
				o.Pool.Logger = mockLogger
				o.TokenProvider = mockProvider
				o.AuditAction = dummyAction
			},
			wantErr: true,
		},
		{
			name: "invalid idle timeout",
			modify: func(o *Options[string]) {
				o.Pool.Name = "valid-pool"
				o.IdleTimeout = -time.Second
				o.Pool.Logger = mockLogger
				o.TokenProvider = mockProvider
				o.AuditAction = dummyAction
			},
			wantErr: true,
		},
		{
			name: "invalid keep alive time",
			modify: func(o *Options[string]) {
				o.Pool.Name = "valid-pool"
				o.KeepAliveTime = 0
				o.Pool.Logger = mockLogger
				o.TokenProvider = mockProvider
				o.AuditAction = dummyAction
			},
			wantErr: true,
		},
		{
			name: "invalid keep alive timeout",
			modify: func(o *Options[string]) {
				o.Pool.Name = "valid-pool"
				o.KeepAliveTimeout = 0
				o.Pool.Logger = mockLogger
				o.TokenProvider = mockProvider
				o.AuditAction = dummyAction
			},
			wantErr: true,
		},
		{
			name: "nil base audit options layer",
			modify: func(o *Options[string]) {
				o.BaseAuditClientOptions = nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewOptions[string]()
			tt.modify(opts)
			err := opts.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
