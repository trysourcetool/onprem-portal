package internal

import (
	"context"

	"github.com/trysourcetool/onprem-portal/internal/core"
)

type ctxKey string

const ContextUserKey ctxKey = "user"

func ContextUser(ctx context.Context) *core.User {
	v, ok := ctx.Value(ContextUserKey).(*core.User)
	if !ok {
		return nil
	}
	return v
}
