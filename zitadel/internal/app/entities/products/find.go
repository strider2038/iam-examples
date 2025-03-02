package products

import (
	"context"

	"app/internal/app/entities/products/domain"
	"app/internal/app/entities/products/view"

	"github.com/muonsoft/errors"
)

type FindQuery struct {
	CompanyID string

	Search string
}

type FindUseCase struct {
	products view.ProductRepository
}

func NewFindUseCase(products view.ProductRepository) *FindUseCase {
	return &FindUseCase{products: products}
}

func (u *FindUseCase) Handle(ctx context.Context, query FindQuery) ([]*view.Product, error) {
	products, err := u.products.Find(ctx, domain.ProductCriteria{
		CompanyID: query.CompanyID,
		Search:    query.Search,
	})
	if err != nil {
		return nil, errors.Errorf("find products: %w", err)
	}

	return products, nil
}
