-- 1. Drop the trigger first (optional but good practice, though dropping the table also drops its triggers)
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- 2. Drop the function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- 3. Drop the table
DROP TABLE IF EXISTS users;