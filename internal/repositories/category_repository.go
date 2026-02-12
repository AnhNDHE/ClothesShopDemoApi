package repositories

import (
	"clothes-shop-api/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	DB *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at, created_by, updated_by, is_active, is_deleted
		FROM categories
		WHERE is_deleted = false
		ORDER BY name
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.ID, &category.Name, &category.Description,
			&category.CreatedAt, &category.UpdatedAt, &category.CreatedBy, &category.UpdatedBy,
			&category.IsActive, &category.IsDeleted,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, name, description string) (*models.Category, error) {
	query := `
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at, created_by, updated_by, is_active, is_deleted
	`

	var category models.Category
	err := r.DB.QueryRow(ctx, query, name, description).Scan(
		&category.ID, &category.Name, &category.Description,
		&category.CreatedAt, &category.UpdatedAt, &category.CreatedBy, &category.UpdatedBy,
		&category.IsActive, &category.IsDeleted,
	)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, id, name, description string) (*models.Category, error) {
	query := `
		UPDATE categories
		SET name = $2, description = $3, updated_at = now()
		WHERE id = $1
		RETURNING id, name, description, created_at, updated_at, created_by, updated_by, is_active, is_deleted
	`

	var category models.Category
	err := r.DB.QueryRow(ctx, query, id, name, description).Scan(
		&category.ID, &category.Name, &category.Description,
		&category.CreatedAt, &category.UpdatedAt, &category.CreatedBy, &category.UpdatedBy,
		&category.IsActive, &category.IsDeleted,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) ToggleCategoryActive(ctx context.Context, id string) error {
	query := `
		UPDATE categories
		SET is_active = NOT is_active, updated_at = now()
		WHERE id = $1
	`

	_, err := r.DB.Exec(ctx, query, id)
	return err
}

func (r *CategoryRepository) SoftDeleteCategory(ctx context.Context, id string) error {
	query := `
		UPDATE categories
		SET is_deleted = true, updated_at = now()
		WHERE id = $1
	`

	_, err := r.DB.Exec(ctx, query, id)
	return err
}
