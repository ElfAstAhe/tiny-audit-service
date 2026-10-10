package grpc

import (
	"testing"
	"time"

	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapAuthDtoSDKToGRPC(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		assert.Nil(t, MapAuthDtoSDKToGRPC(nil))
	})

	t.Run("success mapping", func(t *testing.T) {
		now := time.Now().UTC()
		input := &dto.AuthAuditDTO{
			Source:    "grpc-source",
			EventDate: now,
			Event:     "grpc-login",
			Status:    "success",
		}

		output := MapAuthDtoSDKToGRPC(input)

		require.NotNil(t, output)
		assert.Equal(t, input.Source, output.GetSource())
		assert.Equal(t, input.Event, output.GetEvent())
		assert.Equal(t, now.Unix(), output.GetEventDate().GetSeconds())
	})
}

func TestMapDataValueDTOsSDKToGRPC_FilterNil(t *testing.T) {
	input := []*dto.DataAuditValueDTO{
		{Name: "valid_field_1"},
		nil, // Провоцируем фильтрацию nil
		{Name: "valid_field_2"},
	}

	output := MapDataValueDTOsSDKToGRPC(input)

	require.Len(t, output, 2, "Слайс должен содержать ровно 2 элемента, nil должен отфильтроваться")
	assert.Equal(t, "valid_field_1", output[0].GetName())
	assert.Equal(t, "valid_field_2", output[1].GetName())
}
