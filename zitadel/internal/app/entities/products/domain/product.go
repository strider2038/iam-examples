package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID        uuid.UUID
	CompanyID string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string
}

func NewProduct(name string, companyID string, createdBy string) *Product {
	now := time.Now()

	return &Product{
		ID:        uuid.Must(uuid.NewV7()),
		CompanyID: companyID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: createdBy,
	}
}

type ProductCriteria struct {
	CompanyID string

	Search string
}

type ProductRepository interface {
	Find(ctx context.Context, criteria ProductCriteria) ([]*Product, error)
	Save(ctx context.Context, product *Product) error
}
