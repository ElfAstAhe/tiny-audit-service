package utils

import (
	"errors"
	"testing"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockMarshaller struct {
	mock.Mock
}

func (m *mockMarshaller) Marshal(data any) ([]byte, error) {
	args := m.Called(data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func TestBuildPayload_WithMock(t *testing.T) {
	t.Run("success call", func(t *testing.T) {
		m := new(mockMarshaller)
		input := map[string]string{"key": "val"}
		expectedBytes := []byte(`{"key":"val"}`)

		m.On("Marshal", input).Return(expectedBytes, nil).Once()

		res, err := buildPayload(input, m.Marshal)

		require.NoError(t, err)
		assert.Equal(t, expectedBytes, res)
		m.AssertExpectations(t)
	})

	t.Run("nil data skips marshaller", func(t *testing.T) {
		m := new(mockMarshaller)

		res, err := buildPayload(nil, m.Marshal)

		var invalidArgErr *errs.InvalidArgumentError
		require.Error(t, err)
		assert.ErrorAs(t, err, &invalidArgErr)
		assert.Nil(t, res)
		m.AssertNotCalled(t, "Marshal", mock.Anything)
	})

	t.Run("marshaller error wrapped in CommonError", func(t *testing.T) {
		m := new(mockMarshaller)
		input := "some-data"
		marshalErr := errors.New("encode boom")

		m.On("Marshal", input).Return(nil, marshalErr).Once()

		res, err := buildPayload(input, m.Marshal)

		var commonErr *errs.CommonError
		require.Error(t, err)
		assert.ErrorAs(t, err, &commonErr)
		assert.ErrorIs(t, err, marshalErr)
		assert.Nil(t, res)
		m.AssertExpectations(t)
	})
}
