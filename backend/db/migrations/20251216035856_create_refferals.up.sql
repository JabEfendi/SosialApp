CREATE TABLE IF NOT EXISTS referrals (
    id SERIAL PRIMARY KEY,

    referrer_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    referred_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    status VARCHAR(20) DEFAULT 'pending',
    reward_amount INTEGER DEFAULT 0,
    reward_given_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- ❌ 1 user hanya boleh punya 1 referrer
    CONSTRAINT unique_referred_user UNIQUE (referred_id)
);
