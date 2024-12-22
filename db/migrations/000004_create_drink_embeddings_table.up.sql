CREATE TABLE drink_embeddings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    drink_option_id INTEGER NOT NULL,
    embedding_data TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (drink_option_id) REFERENCES drink_options(id)
); 