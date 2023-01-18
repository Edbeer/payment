DROP TABLE IF EXISTS payment CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE IF NOT EXISTS payment 
(
	payment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    merchant UUID,
    customer UUID,
	card_number VARCHAR(16),
	card_expiry_month VARCHAR(2),
	card_expiry_year VARCHAR(2),
	currency VARCHAR(50),
	operation VARCHAR(50),
	status VARCHAR(50),
	amount serial,
	created_at TIMESTAMP
);
