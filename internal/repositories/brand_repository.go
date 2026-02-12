package repositories

import (
	"clothes-shop-api/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BrandRepository struct {
	DB *pgxpool.Pool
}

func NewBrandRepository(db *pgxpool.Pool) *BrandRepository {
	return &BrandRepository{DB: db}
}

func (r *BrandRepository) GetAllBrands(ctx context.Context) ([]*models.Brand, error) {
	query := `
		SELECT id, name, description, created_at, updated_at, created_by, updated_by, is_active, is_deleted
		FROM brands
		WHERE is_deleted = false
		ORDER BY name
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var brands []*models.Brand
	for rows.Next() {
		var brand models.Brand
		err := rows.Scan(
			&brand.ID, &brand.Name, &brand.Description,
			&brand.CreatedAt, &brand.UpdatedAt, &brand.CreatedBy, &brand.UpdatedBy,
			&brand.IsActive, &brand.IsDeleted,
		)
		if err != nil {
			return nil, err
		}
		brands = append(brands, &brand)
	}

	return brands, nil
}

func (r *BrandRepository) CreateBrand(ctx context.Context, name, description string) (*models.Brand, error) {
	query := `
		INSERT INTO brands (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at, created_by, updated_by, is_active, is_deleted
	`

	var brand models.Brand
	err := r.DB.QueryRow(ctx, query, name, description).Scan(
		&brand.ID, &brand.Name, &brand.Description,
		&brand.CreatedAt, &brand.UpdatedAt, &brand.CreatedBy, &brand.UpdatedBy,
		&brand.IsActive, &brand.IsDeleted,
	)
	if err != nil {
		return nil, err
	}

	return &brand, nil
}

func (r *BrandRepository) UpdateBrand(ctx context.Context, id, name, description string) (*models.Brand, error) {
	query := `
		UPDATE brands
		SET name = $2, description = $3, updated_at = now()
		WHERE id = $1
		RETURNING id, name, description, created_at, updated_at, created_by, updated_by, is_active, is_deleted
	`

	var brand models.Brand
	err := r.DB.QueryRow(ctx, query, id, name, description).Scan(
		&brand.ID, &brand.Name, &brand.Description,
		&brand.CreatedAt, &brand.UpdatedAt, &brand.CreatedBy, &brand.UpdatedBy,
		&brand.IsActive, &brand.IsDeleted,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("brand not found")
		}
		return nil, err
	}

	return &brand, nil
}

func (r *BrandRepository) ToggleBrandActive(ctx context.Context, id string) error {
	query := `
		UPDATE brands
		SET is_active = NOT is_active, updated_at = now()
		WHERE id = $1
	`

	_, err := r.DB.Exec(ctx, query, id)
	return err
}

func (r *BrandRepository) SoftDeleteBrand(ctx context.Context, id string) error {
	query := `
		UPDATE brands
		SET is_deleted = true, updated_at = now()
		WHERE id = $1
	`

	_, err := r.DB.Exec(ctx, query, id)
	return err
}
