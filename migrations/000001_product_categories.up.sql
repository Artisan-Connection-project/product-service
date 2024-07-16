CREATE TABLE product_categories (
    id UUID PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);
SELECT p.id, p.name, p.price, p.category_id 
FROM products p
JOIN product_categories c ON p.category_id = c.id
WHERE c.name LIKE '%Wooden Chair'
AND p.name LIKE '%Wooden'
AND p.price BETWEEN 0.0 AND 100.0
AND p.deleted_at IS NULL
AND c.deleted_at IS NULL;

SELECT * FROM products;
SELECT * FROM product_categories;
SELECT p.id, p.name, p.price, p.category_id 
FROM products p
JOIN product_categories c ON p.category_id = c.id
WHERE c.name LIKE '%home'
AND p.name LIKE '%Wooden'
AND p.price BETWEEN 0.0 AND 100.0
AND p.deleted_at IS NULL
AND c.deleted_at IS NULL;


SELECT p.id, p.name, p.price, p.category_id
FROM products p
JOIN product_categories c ON p.category_id = c.id
WHERE c.name LIKE '%Wooden%'
AND p.name LIKE '%Wooden%'
AND p.price BETWEEN 0.0 AND 100.0
AND p.deleted_at IS NULL
AND c.deleted_at IS NULL;
