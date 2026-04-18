CREATE TABLE IF NOT EXISTS orders (
                                      id SERIAL PRIMARY KEY,
                                      user_id BIGINT NOT NULL REFERENCES users(id),
    product_name VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'new',
    created_at TIMESTAMPTZ DEFAULT NOW()
    );
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);