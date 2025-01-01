-- Create the drink_templates table (template/reference table)
CREATE TABLE IF NOT EXISTS drink_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    size_value REAL NOT NULL CHECK (size_value > 0),
    size_unit TEXT NOT NULL,
    abv REAL NOT NULL CHECK (abv >= 0 AND abv <= 1),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create the drink_log_details table (actual drink information)
CREATE TABLE IF NOT EXISTS drink_log_details (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    size_value REAL NOT NULL CHECK (size_value > 0),
    size_unit TEXT NOT NULL,
    abv REAL NOT NULL CHECK (abv >= 0 AND abv <= 1),
    template_id INTEGER NULL,
    hash_key TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (template_id) REFERENCES drink_templates(id)
);

-- Indexes for drink_templates
CREATE INDEX idx_drink_templates_type ON drink_templates(type);
CREATE INDEX idx_drink_templates_name ON drink_templates(name);

-- Indexes for drink_log_details
CREATE INDEX idx_drink_log_details_type ON drink_log_details(type);
CREATE INDEX idx_drink_log_details_name ON drink_log_details(name);

-- Trigger for drink_templates
CREATE TRIGGER update_drink_templates_timestamp 
AFTER UPDATE ON drink_templates
FOR EACH ROW
BEGIN
    UPDATE drink_templates SET updated_at = CURRENT_TIMESTAMP
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

-- Create drink_logs table with reference to drink_log_details instead of drink_options
CREATE TABLE IF NOT EXISTS drink_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    drink_details_id INTEGER NOT NULL,
    logged_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (drink_details_id) REFERENCES drink_log_details(id)
);

-- Update indexes for drink_logs
CREATE INDEX idx_drink_logs_user_id ON drink_logs(user_id);
CREATE INDEX idx_drink_logs_drink_details_id ON drink_logs(drink_details_id);
CREATE INDEX idx_drink_logs_logged_at ON drink_logs(logged_at);

-- Create drink_embeddings table
CREATE TABLE IF NOT EXISTS drink_embeddings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    drink_template_id INTEGER NOT NULL,
    embedding_data TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (drink_template_id) REFERENCES drink_templates(id)
);

CREATE INDEX idx_drink_embeddings_drink_template_id ON drink_embeddings(drink_template_id); 