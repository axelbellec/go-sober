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
('Wheat Beer', 'beer', 50, 'cl', 0.052),
('Stout', 'beer', 44, 'cl', 0.06),
('Porter', 'beer', 44, 'cl', 0.055),
('Double IPA', 'beer', 33, 'cl', 0.085),
('Pale Ale', 'beer', 33, 'cl', 0.054),
('Belgian Tripel', 'beer', 33, 'cl', 0.085),
('Pilsner', 'beer', 50, 'cl', 0.048),
('Lager', 'beer', 50, 'cl', 0.045),
('Amber Ale', 'beer', 33, 'cl', 0.052),
('Sour Beer', 'beer', 33, 'cl', 0.045),

-- Aperitif
('Aperol Spritz', 'aperitif', 12, 'cl', 0.11),
('Negroni', 'aperitif', 9, 'cl', 0.24),
('Manhattan', 'aperitif', 9, 'cl', 0.25),
('Martini', 'aperitif', 7, 'cl', 0.29),
('Sangria', 'aperitif', 20, 'cl', 0.12),
('Vermouth', 'aperitif', 7, 'cl', 0.18),
('Lillet', 'aperitif', 7, 'cl', 0.17),
('Lillet Rosé', 'aperitif', 7, 'cl', 0.17),
('Lillet Blanc', 'aperitif', 7, 'cl', 0.17),
('Lillet Rouge', 'aperitif', 7, 'cl', 0.17),
('Pastis', 'aperitif', 5, 'cl', 0.45),

-- Wines
('Red Wine', 'wine', 12, 'cl', 0.13),
('White Wine', 'wine', 12, 'cl', 0.115),
('Rosé Wine', 'wine', 12, 'cl', 0.115),
('Champagne', 'wine', 10, 'cl', 0.12),
('Port Wine', 'wine', 6, 'cl', 0.20),
('Prosecco', 'wine', 10, 'cl', 0.11),
('Cava', 'wine', 10, 'cl', 0.115),
('Pinot Noir', 'wine', 12, 'cl', 0.135),
('Cabernet Sauvignon', 'wine', 12, 'cl', 0.14),
('Merlot', 'wine', 12, 'cl', 0.135),
('Chardonnay', 'wine', 12, 'cl', 0.13),
('Sauvignon Blanc', 'wine', 12, 'cl', 0.125),
('Riesling', 'wine', 12, 'cl', 0.11),
('Moscato', 'wine', 12, 'cl', 0.055),
('Sherry', 'wine', 6, 'cl', 0.175),
('Dessert Wine', 'wine', 6, 'cl', 0.16),

-- Cocktails
('Margarita', 'cocktail', 9, 'cl', 0.12),
('Mojito', 'cocktail', 12, 'cl', 0.10),
('Long Island', 'cocktail', 12, 'cl', 0.22),
('Gin & Tonic', 'cocktail', 20, 'cl', 0.08),
('Moscow Mule', 'cocktail', 12, 'cl', 0.10),
('Old Fashioned', 'cocktail', 9, 'cl', 0.32),
('Daiquiri', 'cocktail', 9, 'cl', 0.20),
('Cosmopolitan', 'cocktail', 9, 'cl', 0.20),
('Piña Colada', 'cocktail', 12, 'cl', 0.13),
('Mai Tai', 'cocktail', 12, 'cl', 0.28),
('Espresso Martini', 'cocktail', 9, 'cl', 0.18),
('Dark & Stormy', 'cocktail', 12, 'cl', 0.11),
('Whiskey Sour', 'cocktail', 9, 'cl', 0.28),
('Sex on the Beach', 'cocktail', 12, 'cl', 0.13),
('Cuba Libre', 'cocktail', 12, 'cl', 0.12),

-- Spirits (expanded and updated)
('Whiskey', 'spirit', 4, 'cl', 0.40),
('Scotch Whisky', 'spirit', 4, 'cl', 0.43),
('Bourbon', 'spirit', 4, 'cl', 0.45),
('Irish Whiskey', 'spirit', 4, 'cl', 0.40),
('Rye Whiskey', 'spirit', 4, 'cl', 0.45),
('Vodka', 'spirit', 4, 'cl', 0.40),
('Gin', 'spirit', 4, 'cl', 0.40),
('London Dry Gin', 'spirit', 4, 'cl', 0.45),
('Tequila Blanco', 'spirit', 4, 'cl', 0.38),
('Tequila Reposado', 'spirit', 4, 'cl', 0.40),
('White Rum', 'spirit', 4, 'cl', 0.40),
('Dark Rum', 'spirit', 4, 'cl', 0.40),
('Spiced Rum', 'spirit', 4, 'cl', 0.35),
('Cognac', 'spirit', 4, 'cl', 0.40),
('Brandy', 'spirit', 4, 'cl', 0.35),
('Triple Sec', 'spirit', 4, 'cl', 0.40),
('Kahlua', 'spirit', 4, 'cl', 0.20),
('Baileys', 'spirit', 4, 'cl', 0.17),
('Amaretto', 'spirit', 4, 'cl', 0.28),
('Campari', 'spirit', 4, 'cl', 0.25),
('Vermouth', 'spirit', 4, 'cl', 0.18),

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