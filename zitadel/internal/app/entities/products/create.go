package products

import (
	"context"

	"app/internal/app/entities/products/domain"

	"github.com/google/uuid"
	"github.com/muonsoft/errors"
)

type CreateCommand struct {
	UserID    string
	CompanyID string

	Name string
}

type CreateUseCase struct {
	products domain.ProductRepository
}

func NewCreateUseCase(products domain.ProductRepository) *CreateUseCase {
	return &CreateUseCase{products: products}
}

func (u *CreateUseCase) Handle(ctx context.Context, command CreateCommand) (uuid.UUID, error) {
	product := domain.NewProduct(command.Name, command.CompanyID, command.UserID)
	if err := u.products.Save(ctx, product); err != nil {
		return uuid.Nil, errors.Errorf("save product: %w", err)
	}

	return product.ID, nil
}
