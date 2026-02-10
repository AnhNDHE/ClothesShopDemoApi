CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- USERS
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role TEXT DEFAULT 'customer',
    created_by UUID,
    created_at TIMESTAMP DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMP DEFAULT now(),
    is_active BOOLEAN DEFAULT true,
    is_deleted BOOLEAN DEFAULT false
);

-- CATEGORIES
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    created_by UUID,
    created_at TIMESTAMP DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMP DEFAULT now(),
    is_active BOOLEAN DEFAULT true,
    is_deleted BOOLEAN DEFAULT false
);

-- BRANDS
CREATE TABLE brands (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    created_by UUID,
    created_at TIMESTAMP DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMP DEFAULT now(),
    is_active BOOLEAN DEFAULT true,
    is_deleted BOOLEAN DEFAULT false
);

-- PRODUCTS
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    minprice NUMERIC DEFAULT 0,
    maxprice NUMERIC DEFAULT 0,
    total_stock INT DEFAULT 0,
    category_id UUID REFERENCES categories(id),
    brand_id UUID REFERENCES brands(id),
    created_by UUID,
    created_at TIMESTAMP DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMP DEFAULT now(),
    is_active BOOLEAN DEFAULT true,
    is_deleted BOOLEAN DEFAULT false
);

-- PRODUCT VARIANTS
CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID REFERENCES products(id),
    size TEXT NOT NULL,
    color TEXT NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    price NUMERIC NOT NULL DEFAULT 0,
    image TEXT,
    created_by UUID,
    created_at TIMESTAMP DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMP DEFAULT now(),
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
INSERT INTO products (name, description, minprice, maxprice, total_stock, category_id, brand_id)
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
INSERT INTO products (name, description, minprice, maxprice, total_stock, category_id, brand_id)
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
INSERT INTO products (name, description, minprice, maxprice, total_stock, category_id, brand_id)
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
INSERT INTO products (name, description, minprice, maxprice, total_stock, category_id, brand_id)
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
INSERT INTO products (name, description, minprice, maxprice, total_stock, category_id, brand_id)
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

-- Admin user
INSERT INTO users (email, password, role)
VALUES ('admin@shop.com', '123456', 'admin');
