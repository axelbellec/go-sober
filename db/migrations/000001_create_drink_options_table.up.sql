-- Create the drink_options table
CREATE TABLE drink_options (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    size_value INTEGER NOT NULL,
    size_unit TEXT NOT NULL,
    abv REAL NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert some default drink options
INSERT INTO drink_options (name, type, size_value, size_unit, abv) VALUES
-- Beers
('Beer Pint', 'beer', 50, 'cl', 0.05),
('Beer Half Pint', 'beer', 25, 'cl', 0.05),
('Strong Beer', 'beer', 33, 'cl', 0.075),
('Craft IPA', 'beer', 33, 'cl', 0.065),
('Light Beer', 'beer', 33, 'cl', 0.04),

-- Wines
('Red Wine', 'wine', 12, 'cl', 0.13),
('White Wine', 'wine', 12, 'cl', 0.115),
('Rosé Wine', 'wine', 12, 'cl', 0.115),
('Champagne', 'wine', 10, 'cl', 0.12),
('Port Wine', 'wine', 6, 'cl', 0.20),

-- Cocktails
('Margarita', 'cocktail', 9, 'cl', 0.12),
('Mojito', 'cocktail', 12, 'cl', 0.10),
('Long Island', 'cocktail', 12, 'cl', 0.22),
('Gin & Tonic', 'cocktail', 20, 'cl', 0.08),

-- Spirits
('Whiskey', 'spirit', 4, 'cl', 0.40),
('Vodka', 'spirit', 4, 'cl', 0.40),
('Gin', 'spirit', 4, 'cl', 0.40),
('Tequila', 'spirit', 4, 'cl', 0.38),
('Rum', 'spirit', 4, 'cl', 0.40),
('Liqueur', 'spirit', 4, 'cl', 0.20),

-- Shots
('Tequila Shot', 'shot', 3, 'cl', 0.38),
('Vodka Shot', 'shot', 3, 'cl', 0.40),
('Jägermeister Shot', 'shot', 3, 'cl', 0.35),
('Sambuca Shot', 'shot', 3, 'cl', 0.38),
('Fireball Shot', 'shot', 3, 'cl', 0.33);


CREATE TRIGGER update_drink_options_timestamp 
AFTER UPDATE ON drink_options
FOR EACH ROW
BEGIN
    UPDATE drink_options SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;