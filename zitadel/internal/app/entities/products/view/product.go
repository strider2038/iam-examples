package view

import (
	"context"
	"time"

	"app/internal/app/entities/products/domain"
	users_domain "app/internal/app/entities/users/domain"

	"github.com/google/uuid"
)

type Product struct {
	ID        uuid.UUID
	CompanyID string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy *users_domain.User
}

type ProductRepository interface {
	Find(ctx context.Context, criteria domain.ProductCriteria) ([]*Product, error)
}
