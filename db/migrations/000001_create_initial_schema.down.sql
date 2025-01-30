-- Drop triggers first
DROP TRIGGER IF EXISTS update_drink_log_details_standard_drinks_update;

DROP TRIGGER IF EXISTS update_drink_log_details_standard_drinks_insert;

DROP TRIGGER IF EXISTS update_user_updated_at;

DROP TRIGGER IF EXISTS update_drink_templates_standard_drinks_update;

DROP TRIGGER IF EXISTS update_drink_templates_standard_drinks_insert;

DROP TRIGGER IF EXISTS update_drink_templates_timestamp;

-- Drop indexes
DROP INDEX IF EXISTS idx_drink_logs_logged_at;

DROP INDEX IF EXISTS idx_drink_logs_drink_details_id;

DROP INDEX IF EXISTS idx_drink_logs_user_id;

DROP INDEX IF EXISTS idx_drink_log_details_type;

DROP INDEX IF EXISTS idx_drink_log_details_name;

DROP INDEX IF EXISTS idx_users_email;

DROP INDEX IF EXISTS idx_drink_templates_name;

DROP INDEX IF EXISTS idx_drink_templates_type;

-- Drop tables (in correct order due to foreign key constraints)
DROP TABLE IF EXISTS drink_logs;

DROP TABLE IF EXISTS drink_log_details;

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS drink_templates;