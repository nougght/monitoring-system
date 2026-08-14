package util

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/server/internal/model"
)

func GetAgentIDFromContext(ctx context.Context) (uuid.UUID, error) {
	idFromTLS := ctx.Value(model.ContextKeyAgentID)
	if idFromTLS == nil {
		return uuid.Nil, errors.New("no agent ID in context")
	}
	id, ok := idFromTLS.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("invalid agent ID in context")
	}
	return id, nil
}
