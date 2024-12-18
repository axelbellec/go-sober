-- Create the drink_options table
CREATE TABLE drink_options (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    size_value INTEGER NOT NULL,
    size_unit TEXT NOT NULL,
    abv REAL NOT NULL
);

-- Insert some default drink options
INSERT INTO drink_options (name, type, size_value, size_unit, abv) VALUES
('Beer 50 cl', 'beer', 50, 'cl', 0.05),
('Beer 25 cl', 'beer', 25, 'cl', 0.05),
('Beer 33 cl', 'beer', 33, 'cl', 0.05),
('Wine 12 cl', 'wine', 12, 'cl', 0.12),
('Cocktail 10 cl', 'cocktail', 10, 'cl', 0.12),
('Spirit 4 cl', 'spirit', 4, 'cl', 0.4);

