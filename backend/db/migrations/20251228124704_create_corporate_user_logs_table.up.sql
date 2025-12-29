CREATE TABLE IF NOT EXISTS corporate_user_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    corporate_id BIGINT NOT NULL,
    action VARCHAR(50) NOT NULL,
    created_at TIMESTAMP
);