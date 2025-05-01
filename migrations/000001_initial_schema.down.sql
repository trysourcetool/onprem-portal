BEGIN;

-- First drop all tables
DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS "license";

-- Drop all update_at triggers
DROP TRIGGER IF EXISTS update_user_updated_at ON "user";
DROP TRIGGER IF EXISTS update_license_updated_at ON "license";

END;
