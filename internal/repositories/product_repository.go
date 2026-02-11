package repositories

import (
	"clothes-shop-api/internal/models"
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	DB *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) GetAllProducts(ctx context.Context, page, limit int, minPrice, maxPrice *float64, categoryName, searchName *string) ([]models.Product, error) {
	offset := (page - 1) * limit

	query := `
		SELECT p.id, p.name, p.description, p.min_price, p.max_price, p.total_stock, p.category_id, p.brand_id, p.created_at, p.updated_at, p.is_active, p.is_deleted
		FROM products p
		JOIN categories c ON p.category_id = c.id
		WHERE p.is_active = true AND p.is_deleted = false
	`

	args := []interface{}{}
	argCount := 0

	if minPrice != nil {
		argCount++
		query += ` AND p.min_price >= $` + strconv.Itoa(argCount)
		args = append(args, *minPrice)
	}

	if maxPrice != nil {
		argCount++
		query += ` AND p.max_price <= $` + strconv.Itoa(argCount)
		args = append(args, *maxPrice)
	}

	if categoryName != nil {
		argCount++
		query += ` AND c.name ILIKE $` + strconv.Itoa(argCount)
		args = append(args, "%"+*categoryName+"%")
	}

	if searchName != nil {
		argCount++
		query += ` AND p.name ILIKE $` + strconv.Itoa(argCount)
		args = append(args, "%"+*searchName+"%")
	}

	query += ` ORDER BY p.created_at DESC LIMIT $` + strconv.Itoa(argCount+1) + ` OFFSET $` + strconv.Itoa(argCount+2)
	args = append(args, limit, offset)

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		var brandID *string
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.MinPrice, &product.MaxPrice, &product.TotalStock, &product.CategoryID, &brandID,
			&product.CreatedAt, &product.UpdatedAt, &product.IsActive, &product.IsDeleted)
		if err != nil {
			return nil, err
		}
		if brandID != nil {
			parsedUUID, err := uuid.Parse(*brandID)
			if err != nil {
				return nil, err
			}
			product.BrandID = &parsedUUID
		}
		// Load variants
		variants, err := r.GetProductVariants(ctx, product.ID.String())
		if err != nil {
			return nil, err
		}
		product.Variants = variants
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, name, description string, minPrice, maxPrice float64, totalStock int, categoryName, brandName string) (*models.Product, error) {
	// First, get the category ID by name
	var categoryID string
	err := r.DB.QueryRow(ctx, "SELECT id FROM categories WHERE name = $1", categoryName).Scan(&categoryID)
	if err != nil {
		return nil, err
	}

	// Get brand ID by name if provided
	var brandID *string
	if brandName != "" {
		var bid string
		err := r.DB.QueryRow(ctx, "SELECT id FROM brands WHERE name = $1", brandName).Scan(&bid)
		if err != nil {
			return nil, err
		}
		brandID = &bid
	}

	query := `
		INSERT INTO products (name, description, min_price, max_price, total_stock, category_id, brand_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, description, min_price, max_price, total_stock, category_id, brand_id, created_at, updated_at, is_active, is_deleted, created_by, updated_by
	`

	var product models.Product
	err = r.DB.QueryRow(ctx, query, name, description, minPrice, maxPrice, totalStock, categoryID, brandID).Scan(
		&product.ID, &product.Name, &product.Description, &product.MinPrice, &product.MaxPrice, &product.TotalStock, &product.CategoryID, &brandID,
		&product.CreatedAt, &product.UpdatedAt, &product.IsActive, &product.IsDeleted, &product.CreatedBy, &product.UpdatedBy,
	)
	if err != nil {
		return nil, err
	}

	if brandID != nil {
		parsedUUID, err := uuid.Parse(*brandID)
		if err != nil {
			return nil, err
		}
		product.BrandID = &parsedUUID
	}

	return &product, nil
}

func (r *ProductRepository) GetProductVariants(ctx context.Context, productID string) ([]models.ProductVariant, error) {
	query := `
		SELECT id, product_id, size, color, stock, price, image, created_by, created_at, updated_by, updated_at, is_active, is_deleted
		FROM product_variants
		WHERE product_id = $1 AND is_active = true AND is_deleted = false
		ORDER BY created_at
	`

	rows, err := r.DB.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variants []models.ProductVariant
	for rows.Next() {
		var variant models.ProductVariant
		err := rows.Scan(&variant.ID, &variant.ProductID, &variant.Size, &variant.Color, &variant.Stock, &variant.Price, &variant.Image,
			&variant.CreatedBy, &variant.CreatedAt, &variant.UpdatedBy, &variant.UpdatedAt, &variant.IsActive, &variant.IsDeleted)
		if err != nil {
			return nil, err
		}
		variants = append(variants, variant)
	}

	return variants, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, id string, name, description string, minPrice, maxPrice float64, totalStock int, categoryName, brandName string, updatedBy *string) (*models.Product, error) {
	// First, get the category ID by name
	var categoryID string
	err := r.DB.QueryRow(ctx, "SELECT id FROM categories WHERE name = $1", categoryName).Scan(&categoryID)
	if err != nil {
		return nil, err
	}

	// Get brand ID by name if provided
	var brandID *string
	if brandName != "" {
		var bid string
		err := r.DB.QueryRow(ctx, "SELECT id FROM brands WHERE name = $1", brandName).Scan(&bid)
		if err != nil {
			return nil, err
		}
		brandID = &bid
	}

	query := `
		UPDATE products
		SET name = $2, description = $3, min_price = $4, max_price = $5, total_stock = $6, category_id = $7, brand_id = $8, updated_at = now()
		WHERE id = $1
		RETURNING id, name, description, min_price, max_price, total_stock, category_id, brand_id, created_at, updated_at, is_active
	`

	var product models.Product
	err = r.DB.QueryRow(ctx, query, id, name, description, minPrice, maxPrice, totalStock, categoryID, brandID).Scan(
		&product.ID, &product.Name, &product.Description, &product.MinPrice, &product.MaxPrice, &product.TotalStock, &product.CategoryID, &brandID,
		&product.CreatedAt, &product.UpdatedAt, &product.IsActive,
	)

	if err != nil {
		return nil, err
	}

	if brandID != nil {
		parsedUUID, err := uuid.Parse(*brandID)
		if err != nil {
			return nil, err
		}
		product.BrandID = &parsedUUID
	}

	return &product, nil
}

func (r *ProductRepository) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	query := `
		SELECT id, name
		FROM categories
		ORDER BY name
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *ProductRepository) GetAllBrands(ctx context.Context) ([]models.Brand, error) {
	query := `
		SELECT id, name
		FROM brands
		ORDER BY name
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var brands []models.Brand
	for rows.Next() {
		var brand models.Brand
		err := rows.Scan(&brand.ID, &brand.Name)
		if err != nil {
			return nil, err
		}
		brands = append(brands, brand)
	}

	return brands, nil
}

func (r *ProductRepository) ToggleActive(ctx context.Context, id string, updatedBy *string) (*models.Product, error) {
	query := `
		UPDATE products
		SET is_active = NOT is_active, updated_at = now()
		WHERE id = $1
		RETURNING id, name, description, min_price, max_price, total_stock, category_id, brand_id, created_by, created_at, updated_by, updated_at, is_active, is_deleted
	`

	var product models.Product
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&product.ID, &product.Name, &product.Description, &product.MinPrice, &product.MaxPrice, &product.TotalStock, &product.CategoryID, &product.BrandID,
		&product.CreatedBy, &product.CreatedAt, &product.UpdatedBy, &product.UpdatedAt, &product.IsActive, &product.IsDeleted,
	)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) SoftDelete(ctx context.Context, id string, updatedBy *string) (*models.Product, error) {
	query := `
		UPDATE products
		SET is_deleted = true, updated_at = now()
		WHERE id = $1
		RETURNING id, name, description, min_price, max_price, total_stock, category_id, brand_id, created_by, created_at, updated_by, updated_at, is_active, is_deleted
	`

	var product models.Product
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&product.ID, &product.Name, &product.Description, &product.MinPrice, &product.MaxPrice, &product.TotalStock, &product.CategoryID, &product.BrandID,
		&product.CreatedBy, &product.CreatedAt, &product.UpdatedBy, &product.UpdatedAt, &product.IsActive, &product.IsDeleted,
	)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) CreateProductVariants(ctx context.Context, productID string, variants []struct {
	Size  string
	Color string
	Stock int
	Price float64
	Image string
}) error {
	if len(variants) == 0 {
		return nil
	}

	query := `
		INSERT INTO product_variants (product_id, size, color, stock, price, image)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	for _, variant := range variants {
		_, err := r.DB.Exec(ctx, query, productID, variant.Size, variant.Color, variant.Stock, variant.Price, variant.Image)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ProductRepository) UpdateProductVariants(ctx context.Context, productID string, variants []struct {
	Size  string
	Color string
	Stock int
	Price float64
	Image string
}) error {
	// First, delete existing variants
	_, err := r.DB.Exec(ctx, "DELETE FROM product_variants WHERE product_id = $1", productID)
	if err != nil {
		return err
	}

	// Then insert new variants
	return r.CreateProductVariants(ctx, productID, variants)
}

func (r *ProductRepository) ToggleVariantActive(ctx context.Context, variantID string, updatedBy *string) (*models.ProductVariant, error) {
	query := `
		UPDATE product_variants
		SET is_active = NOT is_active, updated_by = $2, updated_at = now()
		WHERE id = $1
		RETURNING id, product_id, size, color, stock, price, image, created_by, created_at, updated_by, updated_at, is_active, is_deleted
	`

	var variant models.ProductVariant
	var createdByUUID *string
	err := r.DB.QueryRow(ctx, query, variantID, updatedBy).Scan(
		&variant.ID, &variant.ProductID, &variant.Size, &variant.Color, &variant.Stock, &variant.Price, &variant.Image,
		&createdByUUID, &variant.CreatedAt, &variant.UpdatedBy, &variant.UpdatedAt, &variant.IsActive, &variant.IsDeleted,
	)
	if err != nil {
		return nil, err
	}

	return &variant, nil
}

func (r *ProductRepository) SoftDeleteVariant(ctx context.Context, variantID string, updatedBy *string) (*models.ProductVariant, error) {
	query := `
		UPDATE product_variants
		SET is_deleted = true, updated_by = $2, updated_at = now()
		WHERE id = $1
		RETURNING id, product_id, size, color, stock, price, image, created_by, created_at, updated_by, updated_at, is_active, is_deleted
	`

	var variant models.ProductVariant
	var createdByUUID *string
	err := r.DB.QueryRow(ctx, query, variantID, updatedBy).Scan(
		&variant.ID, &variant.ProductID, &variant.Size, &variant.Color, &variant.Stock, &variant.Price, &variant.Image,
		&createdByUUID, &variant.CreatedAt, &variant.UpdatedBy, &variant.UpdatedAt, &variant.IsActive, &variant.IsDeleted,
	)
	if err != nil {
		return nil, err
	}

	return &variant, nil
}
