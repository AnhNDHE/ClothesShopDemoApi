package config

import (
	"context"
	"crypto/tls"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

type Config struct {
	JWTSecret string
}

func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is empty")
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("Cannot parse DB config:", err)
	}

	config.ConnConfig.TLSConfig = &tls.Config{
		InsecureSkipVerify: true, // Required for Render
	}

	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("Cannot connect DB:", err)
	}

	log.Println("Connected to PostgreSQL")
}

func RunMigration() {
	sql := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	-- USERS
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT DEFAULT 'customer',
		created_at TIMESTAMP DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now()
	);

	-- CATEGORIES
	CREATE TABLE IF NOT EXISTS categories (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now()
	);

	-- BRANDS
	CREATE TABLE IF NOT EXISTS brands (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now()
	);

	-- PRODUCTS
	CREATE TABLE IF NOT EXISTS products (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name TEXT NOT NULL,
		description TEXT,
		price NUMERIC NOT NULL,
		stock INT NOT NULL DEFAULT 0,
		category_id UUID REFERENCES categories(id),
		brand_id UUID REFERENCES brands(id),
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now()
	);

	-- PRODUCT VARIANTS
	CREATE TABLE IF NOT EXISTS product_variants (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		product_id UUID REFERENCES products(id) ON DELETE CASCADE,
		size TEXT,
		color TEXT,
		sku TEXT UNIQUE,
		price_adjustment NUMERIC DEFAULT 0,
		stock INT NOT NULL DEFAULT 0,
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now()
	);

	-- CART
	CREATE TABLE IF NOT EXISTS carts (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID REFERENCES users(id),
		created_at TIMESTAMP DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now()
	);

	CREATE TABLE IF NOT EXISTS cart_items (
		cart_id UUID REFERENCES carts(id) ON DELETE CASCADE,
		product_variant_id UUID REFERENCES product_variants(id) ON DELETE CASCADE,
		quantity INT NOT NULL,
		created_at TIMESTAMP DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now(),
		PRIMARY KEY (cart_id, product_variant_id)
	);

	-- ORDERS
	CREATE TABLE IF NOT EXISTS orders (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID REFERENCES users(id),
		total NUMERIC,
		status TEXT DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now()
	);

	-- ORDER ITEMS
	CREATE TABLE IF NOT EXISTS order_items (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
		product_variant_id UUID REFERENCES product_variants(id),
		quantity INT NOT NULL,
		price NUMERIC NOT NULL,
		created_at TIMESTAMP DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now()
	);

	-- SEED DATA
	INSERT INTO categories (name) VALUES
	('T-Shirt'),
	('Jacket'),
	('Pants'),
	('Shoes'),
	('Accessories')
	ON CONFLICT DO NOTHING;

	INSERT INTO brands (name) VALUES
	('Nike'),
	('Adidas'),
	('Puma'),
	('Levi''s'),
	('Zara')
	ON CONFLICT DO NOTHING;

	INSERT INTO products (name, description, price, stock, category_id, brand_id)
	SELECT 'Basic T-Shirt', 'Comfortable cotton t-shirt', 199000, 100, c.id, b.id
	FROM categories c, brands b
	WHERE c.name='T-Shirt' AND b.name='Nike'
	ON CONFLICT DO NOTHING;

	INSERT INTO products (name, description, price, stock, category_id, brand_id)
	SELECT 'Running Jacket', 'Lightweight running jacket', 599000, 50, c.id, b.id
	FROM categories c, brands b
	WHERE c.name='Jacket' AND b.name='Adidas'
	ON CONFLICT DO NOTHING;

	INSERT INTO products (name, description, price, stock, category_id, brand_id)
	SELECT 'Jeans Pants', 'Classic blue jeans', 399000, 75, c.id, b.id
	FROM categories c, brands b
	WHERE c.name='Pants' AND b.name='Levi''s'
	ON CONFLICT DO NOTHING;

	INSERT INTO products (name, description, price, stock, category_id, brand_id)
	SELECT 'Sneakers', 'Comfortable sneakers', 799000, 30, c.id, b.id
	FROM categories c, brands b
	WHERE c.name='Shoes' AND b.name='Puma'
	ON CONFLICT DO NOTHING;

	INSERT INTO products (name, description, price, stock, category_id, brand_id)
	SELECT 'Cap', 'Stylish baseball cap', 99000, 200, c.id, b.id
	FROM categories c, brands b
	WHERE c.name='Accessories' AND b.name='Zara'
	ON CONFLICT DO NOTHING;

	-- Insert variants for each product
	INSERT INTO product_variants (product_id, size, color, sku, stock)
	SELECT p.id, 'S', 'White', 'TSHIRT-S-WHT', 20 FROM products p WHERE p.name='Basic T-Shirt'
	ON CONFLICT DO NOTHING;

	INSERT INTO product_variants (product_id, size, color, sku, stock)
	SELECT p.id, 'M', 'White', 'TSHIRT-M-WHT', 30 FROM products p WHERE p.name='Basic T-Shirt'
	ON CONFLICT DO NOTHING;

	INSERT INTO product_variants (product_id, size, color, sku, stock)
	SELECT p.id, 'L', 'White', 'TSHIRT-L-WHT', 25 FROM products p WHERE p.name='Basic T-Shirt'
	ON CONFLICT DO NOTHING;

	INSERT INTO product_variants (product_id, size, color, sku, stock)
	SELECT p.id, 'M', 'Black', 'JACKET-M-BLK', 25 FROM products p WHERE p.name='Running Jacket'
	ON CONFLICT DO NOTHING;

	INSERT INTO product_variants (product_id, size, color, sku, stock)
	SELECT p.id, 'L', 'Black', 'JACKET-L-BLK', 25 FROM products p WHERE p.name='Running Jacket'
	ON CONFLICT DO NOTHING;

	INSERT INTO product_variants (product_id, size, color, sku, stock)
	SELECT p.id, '30', 'Blue', 'JEANS-30-BLU', 25 FROM products p WHERE p.name='Jeans Pants'
	ON CONFLICT DO NOTHING;

	INSERT INTO product_variants (product_id, size, color, sku, stock)
	SELECT p.id, '32', 'Blue', 'JEANS-32-BLU', 25 FROM products p WHERE p.name='Jeans Pants'
	ON CONFLICT DO NOTHING;

	INSERT INTO product_variants (product_id, size, color, sku, stock)
	SELECT p.id, '34', 'Blue', 'JEANS-34-BLU', 25 FROM products p WHERE p.name='Jeans Pants'
	ON CONFLICT DO NOTHING;

	INSERT INTO product_variants (product_id, size, color, sku, stock)
	SELECT p.id, '8', 'White', 'SNEAKERS-8-WHT', 15 FROM products p WHERE p.name='Sneakers'
	ON CONFLICT DO NOTHING;

	INSERT INTO product_variants (product_id, size, color, sku, stock)
	SELECT p.id, '9', 'White', 'SNEAKERS-9-WHT', 15 FROM products p WHERE p.name='Sneakers'
	ON CONFLICT DO NOTHING;

	INSERT INTO product_variants (product_id, size, color, sku, stock)
	SELECT p.id, 'One Size', 'Black', 'CAP-OS-BLK', 100 FROM products p WHERE p.name='Cap'
	ON CONFLICT DO NOTHING;

	INSERT INTO users (email, password, role)
	VALUES ('admin@shop.com', '$2a$10$examplehashedpassword', 'admin')
	ON CONFLICT DO NOTHING;
	`

	_, err := DB.Exec(context.Background(), sql)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Migration completed successfully")
}

func LoadConfig() Config {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key" // For development, change in production
	}

	return Config{
		JWTSecret: jwtSecret,
	}
}
