DROP TABLE IF EXISTS account CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE IF NOT EXISTS account 
(
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	first_name VARCHAR(50),
	last_name VARCHAR(50),
	card_number VARCHAR(16),
	card_expiry_month VARCHAR(2),
	card_expiry_year VARCHAR(2),
	card_security_code VARCHAR(3),
	balance serial,
	blocked_money serial,
	created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS statement
(
    account_id UUID references account,
    payment_id varchar,
    status varchar
);