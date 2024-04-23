-- Drop the trigger if it exists
DROP TRIGGER IF EXISTS update_users_modtime ON users;

-- Drop the function if it exists
DROP FUNCTION IF EXISTS update_modified_column;

-- Drop the table if it exists
DROP TABLE IF EXISTS users;