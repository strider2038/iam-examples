package repository

import (
	"context"
	"database/sql"
	"strings"

	"app/internal/app/entities/products/domain"

	"github.com/muonsoft/errors"
)

type SQLiteProductRepository struct {
	db *sql.DB
}

func NewSQLiteProductRepository(db *sql.DB) *SQLiteProductRepository {
	return &SQLiteProductRepository{db: db}
}

func (r *SQLiteProductRepository) Find(ctx context.Context, criteria domain.ProductCriteria) ([]*domain.Product, error) {
	query := `SELECT id, company_id, name, created_by, created_at, updated_at FROM products`
	where := make([]string, 0, 2)
	args := make([]any, 0, 2)

	if criteria.CompanyID != "" {
		where = append(where, "company_id = ?")
		args = append(args, criteria.CompanyID)
	}
	if criteria.Search != "" {
		where = append(where, "name LIKE ?")
		args = append(args, "%"+criteria.Search+"%")
	}
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY name ASC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*domain.Product, 0, 10)
	for rows.Next() {
		product := &domain.Product{}
		err := rows.Scan(
			&product.ID,
			&product.CompanyID,
			&product.Name,
			&product.CreatedBy,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Errorf("read products: %w", err)
	}

	return products, nil
}

func (r *SQLiteProductRepository) Save(ctx context.Context, product *domain.Product) error {
	_, err := r.db.Exec(
		"INSERT INTO products (id, company_id, name, created_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		product.ID,
		product.CompanyID,
		product.Name,
		product.CreatedBy,
		product.CreatedAt,
		product.UpdatedAt,
	)
	if err != nil {
		return errors.Errorf("save product: %w", err)
	}

	return nil
}
