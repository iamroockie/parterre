package usecase

import (
	"context"
	"fmt"

	"github.com/iamroockie/parterre/internal/catalog"
)

type CreateHall struct {
	saver hallSaver
}

func NewCreateHall(saver hallSaver) *CreateHall {
	return &CreateHall{saver}
}

func (uc *CreateHall) Execute(
	ctx context.Context,
	p catalog.HallCreateParams,
) (*catalog.Hall, error) {
	hall, err := catalog.NewHall(p)
	if err != nil {
		return nil, fmt.Errorf("create hall: %w", err)
	}

	err = uc.saver.Save(ctx, hall)
	if err != nil {
		return nil, fmt.Errorf("save hall: %w", err)
	}

	return hall, nil
}
