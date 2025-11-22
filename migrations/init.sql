
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    username VARCHAR(255) PRIMARY KEY,
    password VARCHAR(255) NOT NULL,
    balance INTEGER CHECK (balance >= 0)
);

CREATE TABLE items (
    name VARCHAR(255) PRIMARY KEY,
    price INTEGER NOT NULL CHECK (price > 0)
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username_from VARCHAR(255) REFERENCES users(username) ON DELETE CASCADE,
    username_to VARCHAR(255) REFERENCES users(username) ON DELETE CASCADE,
    amount INTEGER NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_items (
    username VARCHAR(255) REFERENCES users(username) ON DELETE CASCADE,
    item_name VARCHAR(255) REFERENCES items(name) ON DELETE CASCADE,
    amount INTEGER CHECK (amount >= 0),
    PRIMARY KEY (username, item_name)
);

INSERT INTO items (name, price) VALUES
('t-shirt', 80),
('cup', 20),
('book', 50),
('pen', 10),
('powerbank', 200),
('hoody', 300),
('umbrella', 200),
('socks', 10),
('wallet', 50),
('pink-hoody', 500);