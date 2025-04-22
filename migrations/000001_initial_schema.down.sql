BEGIN;

-- First drop all tables
DROP TABLE IF EXISTS "user";

-- Drop all update_at triggers
DROP TRIGGER IF EXISTS update_user_updated_at ON "user";

END;
