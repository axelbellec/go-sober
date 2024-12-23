
-- Create the drink_options table
CREATE TABLE IF NOT EXISTS drink_options (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    size_value REAL NOT NULL CHECK (size_value > 0),
    size_unit TEXT NOT NULL,
    abv REAL NOT NULL CHECK (abv >= 0 AND abv <= 1),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for drink_options
CREATE INDEX idx_drink_options_type ON drink_options(type);
CREATE INDEX idx_drink_options_name ON drink_options(name);

-- Trigger for drink_options
CREATE TRIGGER update_drink_options_timestamp 
AFTER UPDATE ON drink_options
FOR EACH ROW
BEGIN
    UPDATE drink_options SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);

CREATE TRIGGER update_user_updated_at
AFTER UPDATE ON users
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Create drink_logs table
CREATE TABLE IF NOT EXISTS drink_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    drink_option_id INTEGER NOT NULL,
    logged_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (drink_option_id) REFERENCES drink_options(id)
);

CREATE INDEX idx_drink_logs_user_id ON drink_logs(user_id);
CREATE INDEX idx_drink_logs_drink_option_id ON drink_logs(drink_option_id);
CREATE INDEX idx_drink_logs_logged_at ON drink_logs(logged_at);

-- Create drink_embeddings table
CREATE TABLE IF NOT EXISTS drink_embeddings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    drink_option_id INTEGER NOT NULL,
    embedding_data TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (drink_option_id) REFERENCES drink_options(id)
);

CREATE INDEX idx_drink_embeddings_drink_option_id ON drink_embeddings(drink_option_id); 