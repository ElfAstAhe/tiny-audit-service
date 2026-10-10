package rest

import (
	"testing"
	"time"

	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapAuthDtoSDKToRest(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		assert.Nil(t, MapAuthDtoSDKToRest(nil))
	})

	t.Run("success mapping", func(t *testing.T) {
		now := time.Now().UTC()
		input := &dto.AuthAuditDTO{
			Source:    "auth-service",
			EventDate: now,
			Event:     "login",
			Status:    "success",
			Username:  "john_doe",
		}

		output := MapAuthDtoSDKToRest(input)

		require.NotNil(t, output)
		assert.Equal(t, input.Source, output.Source)
		assert.Equal(t, now.Format(time.RFC3339), output.EventDate)
		assert.Equal(t, input.Username, output.Username)
	})
}

func TestMapDataValueDTOsSDKToRest_WithNilElements(t *testing.T) {
	input := []*dto.DataAuditValueDTO{
		{Name: "field_1"},
		nil, // Специально подмешиваем nil элемент
		{Name: "field_2"},
	}

	output := MapDataValueDTOsSDKToRest(input)

	// Проверяем, что отфильтровался nil элемент и длина стала равной 2
	require.Len(t, output, 2)
	assert.Equal(t, "field_1", output[0].Name)
	assert.Equal(t, "field_2", output[1].Name)
}
