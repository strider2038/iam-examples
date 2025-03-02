package repository

import (
	"context"

	"app/internal/app/entities/products/domain"
	"app/internal/app/entities/products/view"
	users_domain "app/internal/app/entities/users/domain"
)

type ProductViewRepository struct {
	products domain.ProductRepository
	users    users_domain.UserRepository
}

func NewProductViewRepository(
	products domain.ProductRepository,
	users users_domain.UserRepository,
) *ProductViewRepository {
	return &ProductViewRepository{
		products: products,
		users:    users,
	}
}

func (r *ProductViewRepository) Find(ctx context.Context, criteria domain.ProductCriteria) ([]*view.Product, error) {
	products, err := r.products.Find(ctx, criteria)
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, 0, len(products))
	for _, product := range products {
		userIDs = append(userIDs, product.CreatedBy)
	}

	users, err := r.users.FindByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	out := make([]*view.Product, 0, len(products))
	for _, product := range products {
		out = append(out, &view.Product{
			ID:        product.ID,
			CompanyID: product.CompanyID,
			Name:      product.Name,
			CreatedAt: product.CreatedAt,
			UpdatedAt: product.UpdatedAt,
			CreatedBy: users[product.CreatedBy],
		})
	}

	return out, nil
}
