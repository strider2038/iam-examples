package repository

import (
	"cmp"
	"context"
	"slices"
	"strings"
	"sync"

	"app/internal/app/entities/products/domain"

	"github.com/google/uuid"
)

type InMemoryProductRepository struct {
	mu    sync.Mutex
	items map[uuid.UUID]*domain.Product
}

func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		items: make(map[uuid.UUID]*domain.Product),
	}
}

func (r *InMemoryProductRepository) Find(ctx context.Context, criteria domain.ProductCriteria) ([]*domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	products := make([]*domain.Product, 0, len(r.items))
	for _, product := range r.items {
		if matches(product, criteria) {
			products = append(products, product)
		}
	}

	slices.SortFunc(products, func(a, b *domain.Product) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return products, nil
}

func (r *InMemoryProductRepository) Save(ctx context.Context, product *domain.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.items[product.ID] = product

	return nil
}

func matches(product *domain.Product, criteria domain.ProductCriteria) bool {
	if criteria.CompanyID != "" && product.CompanyID != criteria.CompanyID {
		return false
	}

	if criteria.Search != "" && !strings.Contains(strings.ToLower(product.Name), strings.ToLower(criteria.Search)) {
		return false
	}

	return true
}
