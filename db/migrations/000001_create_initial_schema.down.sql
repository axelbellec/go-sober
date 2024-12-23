
-- Drop indexes
DROP INDEX IF EXISTS idx_drink_embeddings_drink_option_id;
DROP INDEX IF EXISTS idx_drink_logs_logged_at;
DROP INDEX IF EXISTS idx_drink_logs_drink_option_id;
DROP INDEX IF EXISTS idx_drink_logs_user_id;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_drink_options_name;
DROP INDEX IF EXISTS idx_drink_options_type;

-- Drop tables (in correct order due to foreign key constraints)
DROP TABLE IF EXISTS drink_embeddings;
DROP TABLE IF EXISTS drink_logs;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS drink_options; 