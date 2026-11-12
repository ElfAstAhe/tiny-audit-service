package utils

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/ext/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/http"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultBuildDataAuditOutbox(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	var typedNilBuilder *DataAuditBuilder

	t.Run("success build outbox entity", func(t *testing.T) {
		builder := NewDataAuditBuilder().
			WithSource("tiny-auth-service").
			WithEventDate(now).
			WithInstanceID("user-123").
			WithEvent(dto.DataEventCreate)

		res, err := DefaultBuildDataAuditOutbox(context.Background(), builder)

		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "user-123", res.InstanceID)
		assert.Equal(t, http.MediaTypeApplicationJSON, res.ContentType)
		assert.Equal(t, domain.OutboxStatusPending, res.Status)
		assert.Equal(t, now, res.CreatedAt)

		var parsedDTO dto.DataAuditDTO
		err = json.Unmarshal(res.Payload, &parsedDTO)
		require.NoError(t, err)
		assert.Equal(t, "tiny-auth-service", parsedDTO.Source)
		assert.Equal(t, "user-123", parsedDTO.InstanceID)
	})

	t.Run("nil builder returns InvalidArgumentError", func(t *testing.T) {
		res, err := DefaultBuildDataAuditOutbox(context.Background(), nil)

		var invalidArgErr *errs.InvalidArgumentError
		require.Error(t, err)
		assert.ErrorAs(t, err, &invalidArgErr)
		assert.Nil(t, res)
	})

	t.Run("typed nil builder returns InvalidArgumentError", func(t *testing.T) {
		res, err := DefaultBuildDataAuditOutbox(context.Background(), typedNilBuilder)

		var invalidArgErr *errs.InvalidArgumentError
		require.Error(t, err)
		assert.ErrorAs(t, err, &invalidArgErr)
		assert.Nil(t, res)
	})
}
