CREATE TABLE topup_transactions (
    id BIGSERIAL PRIMARY KEY,
    invoice_number UUID NOT NULL UNIQUE,
    user_id BIGINT NOT NULL REFERENCES users(id),
    token_package_id BIGINT REFERENCES token_packages(id),
    token_amount BIGINT NOT NULL,
    price BIGINT NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    payment_reference VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);