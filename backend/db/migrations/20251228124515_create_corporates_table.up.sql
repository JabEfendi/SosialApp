CREATE TABLE IF NOT EXISTS corporates (
    id BIGSERIAL PRIMARY KEY,

    -- corporate identity
    name VARCHAR(150) NOT NULL,
    logo VARCHAR(255),
    reffcorporate VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    status VARCHAR(50) DEFAULT 'active',

    -- corporate contact (NON LOGIN)
    email_corporate VARCHAR(100),
    phone VARCHAR(20),
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100),
    zip_code VARCHAR(20),

    -- PIC (LOGIN ACCOUNT)
    email VARCHAR(100) NOT NULL UNIQUE,     -- email PIC
    password VARCHAR(255) NOT NULL,          -- password PIC
    name_PIC VARCHAR(150),
    phone_PIC VARCHAR(20),
    age_PIC INT,

    two_fa_enabled BOOLEAN DEFAULT FALSE,
    two_fa_secret VARCHAR(100),

    created_by BIGINT NULL,
    updated_by BIGINT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
