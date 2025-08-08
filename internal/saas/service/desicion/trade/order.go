package trade

import (
    "context"

    "study/internal/saas/model"
)

type OrderHandler interface {
    HandleOrder(ctx context.Context, order model.Order) (err error)
}
