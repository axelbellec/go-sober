-- Create the drink_templates table (template/reference table)
CREATE TABLE
    IF NOT EXISTS drink_templates (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        type TEXT NOT NULL,
        size_value REAL NOT NULL CHECK (size_value > 0),
        size_unit TEXT NOT NULL,
        abv REAL NOT NULL CHECK (
            abv >= 0
            AND abv <= 1
        ),
        standard_drinks REAL DEFAULT 0 CHECK (standard_drinks >= 0),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT NULL
    );

-- Indexes for drink_templates
CREATE INDEX idx_drink_templates_type ON drink_templates (type);

CREATE INDEX idx_drink_templates_name ON drink_templates (name);

-- Split the trigger into two separate triggers for INSERT and UPDATE
DROP TRIGGER IF EXISTS update_drink_templates_timestamp;

DROP TRIGGER IF EXISTS update_drink_templates_standard_drinks_insert;

DROP TRIGGER IF EXISTS update_drink_templates_standard_drinks_update;

-- Trigger for drink_templates
CREATE TRIGGER update_drink_templates_timestamp AFTER
UPDATE ON drink_templates FOR EACH ROW BEGIN
UPDATE drink_templates
SET
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = OLD.id;

END;

-- Add triggers for drink_templates standard_drinks calculation
CREATE TRIGGER update_drink_templates_standard_drinks_insert AFTER INSERT ON drink_templates FOR EACH ROW BEGIN
UPDATE drink_templates
SET
    abv = round(NEW.abv, 3),
    standard_drinks = round(
        CASE
            WHEN NEW.size_unit = 'ml' THEN (NEW.size_value * NEW.abv * 0.789) / 10
            WHEN NEW.size_unit = 'cl' THEN (NEW.size_value * 10 * NEW.abv * 0.789) / 10
            ELSE 0
        END,
        4
    )
WHERE
    id = NEW.id;

END;

CREATE TRIGGER update_drink_templates_standard_drinks_update AFTER
UPDATE ON drink_templates FOR EACH ROW BEGIN
UPDATE drink_templates
SET
    abv = round(NEW.abv, 3),
    standard_drinks = round(
        CASE
            WHEN NEW.size_unit = 'ml' THEN (NEW.size_value * NEW.abv * 0.789) / 10
            WHEN NEW.size_unit = 'cl' THEN (NEW.size_value * 10 * NEW.abv * 0.789) / 10
            ELSE 0
        END,
        4
    )
WHERE
    id = NEW.id;

END;

-- Create the drink_log_details table (actual drink information)
CREATE TABLE
    IF NOT EXISTS drink_log_details (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        type TEXT NOT NULL,
        size_value REAL NOT NULL CHECK (size_value > 0),
        size_unit TEXT NOT NULL,
        abv REAL NOT NULL CHECK (
            abv >= 0
            AND abv <= 1
        ),
        standard_drinks REAL DEFAULT 0 CHECK (standard_drinks >= 0),
        template_id INTEGER NULL,
        hash_key TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT NULL,
        FOREIGN KEY (template_id) REFERENCES drink_templates (id)
    );

-- Indexes for drink_log_details
CREATE INDEX idx_drink_log_details_type ON drink_log_details (type);

CREATE INDEX idx_drink_log_details_name ON drink_log_details (name);

-- Split the trigger into two separate triggers for INSERT and UPDATE
DROP TRIGGER IF EXISTS update_drink_log_details_standard_drinks_insert;

DROP TRIGGER IF EXISTS update_drink_log_details_standard_drinks_update;

CREATE TRIGGER update_drink_log_details_standard_drinks_insert AFTER INSERT ON drink_log_details FOR EACH ROW BEGIN
UPDATE drink_log_details
SET
    abv = round(NEW.abv, 3),
    standard_drinks = round(
        CASE
            WHEN NEW.size_unit = 'ml' THEN (NEW.size_value * NEW.abv * 0.789) / 10
            WHEN NEW.size_unit = 'cl' THEN (NEW.size_value * 10 * NEW.abv * 0.789) / 10
            ELSE 0
        END,
        4
    )
WHERE
    id = NEW.id;

END;

CREATE TRIGGER update_drink_log_details_standard_drinks_update AFTER
UPDATE ON drink_log_details FOR EACH ROW BEGIN
UPDATE drink_log_details
SET
    abv = round(NEW.abv, 3),
    standard_drinks = round(
        CASE
            WHEN NEW.size_unit = 'ml' THEN (NEW.size_value * NEW.abv * 0.789) / 10
            WHEN NEW.size_unit = 'cl' THEN (NEW.size_value * 10 * NEW.abv * 0.789) / 10
            ELSE 0
        END,
        4
    )
WHERE
    id = NEW.id;

END;

-- Create users table
CREATE TABLE
    IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT NULL
    );

CREATE INDEX idx_users_email ON users (email);

CREATE TRIGGER update_user_updated_at AFTER
UPDATE ON users BEGIN
UPDATE users
SET
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = NEW.id;

END;

-- Create drink_logs table with reference to drink_log_details instead of drink_options
CREATE TABLE
    IF NOT EXISTS drink_logs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        drink_details_id INTEGER NOT NULL,
        logged_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT NULL,
        FOREIGN KEY (user_id) REFERENCES users (id),
        FOREIGN KEY (drink_details_id) REFERENCES drink_log_details (id)
    );

-- Update indexes for drink_logs
CREATE INDEX idx_drink_logs_user_id ON drink_logs (user_id);

CREATE INDEX idx_drink_logs_drink_details_id ON drink_logs (drink_details_id);

CREATE INDEX idx_drink_logs_logged_at ON drink_logs (logged_at);

---
CREATE TABLE
    IF NOT EXISTS user_profiles (
        user_id INTEGER PRIMARY KEY,
        weight_kg REAL NOT NULL CHECK (weight_kg > 0),
        gender TEXT NOT NULL CHECK (gender IN ('male', 'female', 'unknown')),
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users (id)
    );

CREATE INDEX idx_user_profiles_user_id ON user_profiles (user_id);