CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    product_name VARCHAR(100) NOT NULL,
    quantity INT NOT NULL,
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    attachment_key VARCHAR(255) DEFAULT NULL
);
