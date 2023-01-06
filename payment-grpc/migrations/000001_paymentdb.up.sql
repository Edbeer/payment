DROP TABLE IF EXISTS payment CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE IF NOT EXISTS payment 
(
	payment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    payment_receiver UUID,
    payer UUID,
	currency VARCHAR(50),
	operation VARCHAR(50),
	status VARCHAR(50),
	amount serial,
	created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS statement
(
    payment_id UUID references payment,
    status varchar
);