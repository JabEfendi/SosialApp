CREATE TABLE account_banks (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    logo VARCHAR(255) NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL
);

INSERT INTO account_banks (code, name, type, logo, is_active, created_at, updated_at) VALUES
-- BANK
('BCA', 'Bank Central Asia', 'bank', 'bca.png', true, NOW(), NOW()),
('BRI', 'Bank Rakyat Indonesia', 'bank', 'bri.png', true, NOW(), NOW()),
('BNI', 'Bank Negara Indonesia', 'bank', 'bni.png', true, NOW(), NOW()),
('MANDIRI', 'Bank Mandiri', 'bank', 'mandiri.png', true, NOW(), NOW()),
('CIMB', 'CIMB Niaga', 'bank', 'cimb.png', true, NOW(), NOW()),
('PERMATA', 'Bank Permata', 'bank', 'permata.png', true, NOW(), NOW()),
('DANAMON', 'Bank Danamon', 'bank', 'danamon.png', true, NOW(), NOW()),
('BTN', 'Bank Tabungan Negara', 'bank', 'btn.png', true, NOW(), NOW()),

-- E-WALLET
('OVO', 'OVO', 'ewallet', 'ovo.png', true, NOW(), NOW()),
('GOPAY', 'GoPay', 'ewallet', 'gopay.png', true, NOW(), NOW()),
('DANA', 'DANA', 'ewallet', 'dana.png', true, NOW(), NOW()),
('SHOPEEPAY', 'ShopeePay', 'ewallet', 'shopeepay.png', true, NOW(), NOW()),
('LINKAJA', 'LinkAja', 'ewallet', 'linkaja.png', true, NOW(), NOW());