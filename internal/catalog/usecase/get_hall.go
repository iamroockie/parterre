package usecase

import (
	"context"
	"uuid"

	"github.com/iamroockie/parterre/internal/catalog"
)

type GetHall struct {
	getter hallGetter
}

func NewGetHall(getter hallGetter) *GetHall {
	return &GetHall{getter}
}

func (uc GetHall) Execute(ctx context.Context, id uuid.UUID) (*catalog.Hall, error) {
	return uc.getter.Get(ctx, id)
}
