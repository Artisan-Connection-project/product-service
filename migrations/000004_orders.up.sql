CREATE TYPE order_status AS ENUM (
    'pending',
    'processing',
    'shipped',
    'delivered',
    'canceled',
    'returned'
);

CREATE TABLE orders (
    id UUID PRIMARY KEY,
    user_id UUID,
    total_amount DECIMAL(10, 2) NOT NULL,
    status order_status NOT NULL,
    shipping_address JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);
ALTER TABLE orders
ADD COLUMN tracking_number VARCHAR(255),
ADD COLUMN carrier VARCHAR(255),
ADD COLUMN estimated_delivery_date TIMESTAMP;
