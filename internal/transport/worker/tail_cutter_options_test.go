package worker

import (
	"context"
	"testing"
	"time"

	mocks2 "github.com/ElfAstAhe/go-service-template/pkg/logger/mocks"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTailCutterOptions_Defaults(t *testing.T) {
	opts := NewTailCutterOptions()

	require.NotNil(t, opts.BaseSchedulerDispatcherOptions)
	assert.Equal(t, 1, opts.WorkerCount)
	assert.Equal(t, 32, opts.DataCapacity)
	assert.False(t, opts.CompleteProcess)
	assert.Equal(t, 30*time.Second, opts.StartInterval)
	assert.Equal(t, 5*time.Minute, opts.ScheduleInterval)
	assert.Equal(t, 14*24*time.Hour, opts.DataInterval)
	assert.False(t, opts.CutEnabled)
}

func TestTailCutterOptions_WithFunctions(t *testing.T) {
	opts := NewTailCutterOptions()

	mockGetUC := mocks.NewMockTailGetUseCase[string](t)
	mockCutUC := mocks.NewMockTailCutUseCase[string](t)

	WithTailCutName("custom-cutter")(opts)
	WithTailCutWorkerCount(5)(opts)
	WithTailCutDataInterval(24 * time.Hour)(opts)
	WithTailCutCutEnabled(true)(opts)
	WithTailCutTailGetUC(mockGetUC)(opts)
	WithTailCutterTailCutUC(mockCutUC)(opts)

	assert.Equal(t, "custom-cutter", opts.Name)
	assert.Equal(t, 5, opts.WorkerCount)
	assert.Equal(t, 24*time.Hour, opts.DataInterval)
	assert.True(t, opts.CutEnabled)
	assert.Equal(t, mockGetUC, opts.TailGetUC)
	assert.Equal(t, mockCutUC, opts.TailCutUC)
}

func TestTailCutterOptions_Validate(t *testing.T) {
	mockLogger := mocks2.NewMockLogger(t)
	var mockJobHandler worker.JobHandler[string] = func(ctx context.Context, workerIndex int, data string) error { return nil }
	var mockDataProvider worker.DispatcherDataProvider[string] = func(ctx context.Context, event time.Time) ([]string, error) { return nil, nil }
	mockGetUC := mocks.NewMockTailGetUseCase[string](t)
	mockCutUC := mocks.NewMockTailCutUseCase[string](t)

	validOpts := func() *TailCutterOptions {
		opts := NewTailCutterOptions()
		opts.Name = "valid-tail-cutter"
		opts.StopTimeout = time.Second * 5
		opts.JobHandler = mockJobHandler
		opts.Logger = mockLogger
		opts.DataProvider = mockDataProvider
		WithTailCutTailGetUC(mockGetUC)(opts)
		WithTailCutterTailCutUC(mockCutUC)(opts)
		return opts
	}

	tests := []struct {
		name        string
		modify      func(opts *TailCutterOptions)
		wantErrText string
	}{
		{
			name:        "Success valid options",
			modify:      func(opts *TailCutterOptions) {},
			wantErrText: "",
		},
		{
			name: "Error empty name",
			modify: func(opts *TailCutterOptions) {
				opts.Name = ""
			},
			wantErrText: "name is required",
		},
		{
			name: "Error data interval zero",
			modify: func(opts *TailCutterOptions) {
				opts.DataInterval = 0
			},
			wantErrText: "data interval is required",
		},
		{
			name: "Error nil tail get usecase",
			modify: func(opts *TailCutterOptions) {
				WithTailCutTailGetUC(nil)(opts)
			},
			wantErrText: "tail get use case is required",
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
