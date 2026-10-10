package grpc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRawClientOptions_Defaults(t *testing.T) {
	opts := NewRawClientOptions()

	assert.Equal(t, DefaultSecure, opts.Secure)
	assert.Equal(t, DefaultConnectionTimeout, opts.ConnTimeout)
	assert.Equal(t, DefaultKeepAliveTime, opts.KATime)
	assert.Equal(t, DefaultKeepAliveTimeout, opts.KATimeout)
	assert.Equal(t, DefaultKeepAlivePermitWithoutStream, opts.KAPermitWithoutStream)
}

func TestRawClientOptions_WithFunctions(t *testing.T) {
	opts := NewRawClientOptions()

	WithRawSecure(true)(opts)
	WithRawConnectionTimeout(5 * time.Second)(opts)
	WithRawKeepAliveTime(30 * time.Second)(opts)
	WithRawKeepAliveTimeout(10 * time.Second)(opts)
	WithRawKeepAlivePermitWithoutStream(false)(opts)

	assert.True(t, opts.Secure)
	assert.Equal(t, 5*time.Second, opts.ConnTimeout)
	assert.Equal(t, 30*time.Second, opts.KATime)
	assert.Equal(t, 10*time.Second, opts.KATimeout)
	assert.False(t, opts.KAPermitWithoutStream)
}

//goland:noinspection DuplicatedCode
func TestRawClientOptions_Validate(t *testing.T) {
	validOpts := func() *RawClientOptions {
		return NewRawClientOptions()
	}

	tests := []struct {
		name        string
		modify      func(opts *RawClientOptions)
		wantErrText string
	}{
		{
			name:        "Success valid options",
			modify:      func(opts *RawClientOptions) {},
			wantErrText: "",
		},
		{
			name: "Error zero connection timeout",
			modify: func(opts *RawClientOptions) {
				opts.ConnTimeout = 0
			},
			wantErrText: "connection timeout is required",
		},
		{
			name: "Error negative ka time",
			modify: func(opts *RawClientOptions) {
				opts.KATime = -1 * time.Second
			},
			wantErrText: "ka time is required",
		},
		{
			name: "Error zero ka timeout",
			modify: func(opts *RawClientOptions) {
				opts.KATimeout = 0
			},
			wantErrText: "ka timeout is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := validOpts()
			tt.modify(opts)

			err := opts.Validate()

			if tt.wantErrText != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrText)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
