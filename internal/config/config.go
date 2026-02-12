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
	JWTSecret        string
	SMTPHost         string
	SMTPPort         string
	SMTPUsername     string
	SMTPPassword     string
	EmailFrom        string
	AppBaseURLLocal  string
	AppBaseURLDeploy string
}

// InitDB initializes the PostgreSQL connection
func InitDB() {
	dsn := os.Getenv("DATABASE_URL")

	// If DATABASE_URL is not set, build from individual env vars (for local development or Render)
	if dsn == "" {
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbName := os.Getenv("DB_NAME")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbSSLMode := os.Getenv("DB_SSLMODE")

		if dbHost == "" || dbPort == "" || dbName == "" || dbUser == "" || dbPassword == "" {
			log.Fatal("DATABASE_URL or individual DB_* environment variables are required")
		}

		if dbSSLMode == "" {
			dbSSLMode = "disable"
		}

		dsn = "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=" + dbSSLMode
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("Cannot parse DB config:", err)
	}

	// Set TLS config for external databases (like Render) when sslmode=require
	if os.Getenv("DB_SSLMODE") == "require" || os.Getenv("DATABASE_URL") != "" {
		config.ConnConfig.TLSConfig = &tls.Config{
			InsecureSkipVerify: true, // Required for Render
		}
	}

	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("Cannot connect DB:", err)
	}

	log.Println("Connected to PostgreSQL")
}

func RunMigration() {
	sql := `
	-- Drop tables if they exist (for development)
	DROP TABLE IF EXISTS orders CASCADE;
	DROP TABLE IF EXISTS cart_items CASCADE;
	DROP TABLE IF EXISTS carts CASCADE;
	DROP TABLE IF EXISTS product_variants CASCADE;
	DROP TABLE IF EXISTS products CASCADE;
	DROP TABLE IF EXISTS brands CASCADE;
	DROP TABLE IF EXISTS categories CASCADE;
	DROP TABLE IF EXISTS users CASCADE;

	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- USERS
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role TEXT DEFAULT 'customer',
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    created_by UUID,
    updated_by UUID,
    is_active BOOLEAN DEFAULT false,
    is_deleted BOOLEAN DEFAULT false
);

-- CATEGORIES
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    created_by UUID,
    updated_by UUID,
    is_active BOOLEAN DEFAULT true,
    is_deleted BOOLEAN DEFAULT false
);

-- BRANDS
CREATE TABLE brands (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    created_by UUID,
    updated_by UUID,
    is_active BOOLEAN DEFAULT true,
    is_deleted BOOLEAN DEFAULT false
);

-- PRODUCTS
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    min_price NUMERIC NOT NULL,
    max_price NUMERIC NOT NULL,
    total_stock INT NOT NULL,
    category_id UUID REFERENCES categories(id),
    brand_id UUID REFERENCES brands(id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    created_by UUID,
    updated_by UUID,
    is_active BOOLEAN DEFAULT true,
    is_deleted BOOLEAN DEFAULT false
);

-- PRODUCT_VARIANTS
CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID REFERENCES products(id),
    size TEXT,
    color TEXT,
    stock INT NOT NULL,
    price NUMERIC NOT NULL,
    image TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    created_by UUID,
    updated_by UUID,
    is_active BOOLEAN DEFAULT true,
    is_deleted BOOLEAN DEFAULT false
);

-- CART
CREATE TABLE carts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id)
);

CREATE TABLE cart_items (
    cart_id UUID REFERENCES carts(id),
    product_id UUID REFERENCES products(id),
    quantity INT NOT NULL,
    PRIMARY KEY (cart_id, product_id)
);

-- ORDERS
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    total NUMERIC,
    status TEXT,
    created_at TIMESTAMP DEFAULT now()
);

-- SEED DATA

-- Admin user (password is hashed version of '123456')
INSERT INTO users (email, password, role)
VALUES ('admin@shop.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin');

-- Categories (5 records)
INSERT INTO categories (name, description) VALUES
('T-Shirt', 'Comfortable cotton t-shirts'),
('Jacket', 'Stylish jackets for all seasons'),
('Pants', 'Various types of pants and trousers'),
('Shoes', 'Footwear for men and women'),
('Accessories', 'Fashion accessories and jewelry');

-- Brands (5 records)
INSERT INTO brands (name, description) VALUES
('Nike', 'Leading sportswear brand'),
('Adidas', 'Global sports and lifestyle brand'),
('Zara', 'Fast fashion retailer'),
('H&M', 'Affordable fashion brand'),
('Levi''s', 'Iconic denim brand');

-- Products (5 records, each with 2 variants)
-- Product 1: Nike T-Shirt
INSERT INTO products (name, description, min_price, max_price, total_stock, category_id, brand_id)
SELECT 'Nike Sport T-Shirt', 'Comfortable athletic t-shirt', 250000, 300000, 150,
       c.id, b.id
FROM categories c, brands b
WHERE c.name = 'T-Shirt' AND b.name = 'Nike';

-- Variants for Nike T-Shirt
INSERT INTO product_variants (product_id, size, color, stock, price, image)
SELECT p.id, 'M', 'Black', 50, 250000, 'nike-tshirt-black-m.jpg'
FROM products p WHERE p.name = 'Nike Sport T-Shirt';

INSERT INTO product_variants (product_id, size, color, stock, price, image)
SELECT p.id, 'L', 'White', 100, 300000, 'nike-tshirt-white-l.jpg'
FROM products p WHERE p.name = 'Nike Sport T-Shirt';

-- Product 2: Adidas Jacket
INSERT INTO products (name, description, min_price, max_price, total_stock, category_id, brand_id)
SELECT 'Adidas Winter Jacket', 'Warm winter jacket with hood', 800000, 900000, 80,
       c.id, b.id
FROM categories c, brands b
WHERE c.name = 'Jacket' AND b.name = 'Adidas';

-- Variants for Adidas Jacket
INSERT INTO product_variants (product_id, size, color, stock, price, image)
SELECT p.id, 'M', 'Blue', 40, 800000, 'adidas-jacket-blue-m.jpg'
FROM products p WHERE p.name = 'Adidas Winter Jacket';

INSERT INTO product_variants (product_id, size, color, stock, price, image)
SELECT p.id, 'L', 'Black', 40, 900000, 'adidas-jacket-black-l.jpg'
FROM products p WHERE p.name = 'Adidas Winter Jacket';

-- Product 3: Zara Pants
INSERT INTO products (name, description, min_price, max_price, total_stock, category_id, brand_id)
SELECT 'Zara Slim Fit Pants', 'Elegant slim fit trousers', 450000, 500000, 120,
       c.id, b.id
FROM categories c, brands b
WHERE c.name = 'Pants' AND b.name = 'Zara';

-- Variants for Zara Pants
INSERT INTO product_variants (product_id, size, color, stock, price, image)
SELECT p.id, '32', 'Gray', 60, 450000, 'zara-pants-gray-32.jpg'
FROM products p WHERE p.name = 'Zara Slim Fit Pants';

INSERT INTO product_variants (product_id, size, color, stock, price, image)
SELECT p.id, '34', 'Black', 60, 500000, 'zara-pants-black-34.jpg'
FROM products p WHERE p.name = 'Zara Slim Fit Pants';

-- Product 4: H&M Shoes
INSERT INTO products (name, description, min_price, max_price, total_stock, category_id, brand_id)
SELECT 'H&M Casual Sneakers', 'Comfortable everyday sneakers', 350000, 400000, 200,
       c.id, b.id
FROM categories c, brands b
WHERE c.name = 'Shoes' AND b.name = 'H&M';

-- Variants for H&M Shoes
INSERT INTO product_variants (product_id, size, color, stock, price, image)
SELECT p.id, '42', 'White', 100, 350000, 'hm-sneakers-white-42.jpg'
FROM products p WHERE p.name = 'H&M Casual Sneakers';

INSERT INTO product_variants (product_id, size, color, stock, price, image)
SELECT p.id, '43', 'Blue', 100, 400000, 'hm-sneakers-blue-43.jpg'
FROM products p WHERE p.name = 'H&M Casual Sneakers';

-- Product 5: Levi's Accessories (Belt)
INSERT INTO products (name, description, min_price, max_price, total_stock, category_id, brand_id)
SELECT 'Levi''s Leather Belt', 'Classic leather belt', 150000, 180000, 90,
       c.id, b.id
FROM categories c, brands b
WHERE c.name = 'Accessories' AND b.name = 'Levi''s';

-- Variants for Levi's Belt
INSERT INTO product_variants (product_id, size, color, stock, price, image)
SELECT p.id, 'M', 'Brown', 45, 150000, 'levis-belt-brown-m.jpg'
FROM products p WHERE p.name = 'Levi''s Leather Belt';

INSERT INTO product_variants (product_id, size, color, stock, price, image)
SELECT p.id, 'L', 'Black', 45, 180000, 'levis-belt-black-l.jpg'
FROM products p WHERE p.name = 'Levi''s Leather Belt';

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

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	emailFrom := os.Getenv("EMAIL_FROM")
	appBaseURLLocal := os.Getenv("APP_BASE_URL_LOCAL")
	appBaseURLDeploy := os.Getenv("APP_BASE_URL_DEPLOY")

	return Config{
		JWTSecret:        jwtSecret,
		SMTPHost:         smtpHost,
		SMTPPort:         smtpPort,
		SMTPUsername:     smtpUsername,
		SMTPPassword:     smtpPassword,
		EmailFrom:        emailFrom,
		AppBaseURLLocal:  appBaseURLLocal,
		AppBaseURLDeploy: appBaseURLDeploy,
	}
}
