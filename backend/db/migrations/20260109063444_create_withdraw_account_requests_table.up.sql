CREATE TABLE withdraw_account_requests (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    withdraw_account_id INT NOT NULL REFERENCES user_withdraw_accounts(id) ON DELETE CASCADE,
    account_bank_id INT NOT NULL REFERENCES account_banks(id) ON DELETE CASCADE,
    account_number VARCHAR(50) NOT NULL,
    account_name VARCHAR(100) NOT NULL,
    status VARCHAR DEFAULT 'pending',
    auto_approve_at TIMESTAMP NOT NULL,
    approved_at TIMESTAMP NULL,
    rejected_reason TEXT NULL,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL
);