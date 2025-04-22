BEGIN;

-- Function to automatically update updated_at columns
CREATE OR REPLACE FUNCTION update_updated_at_column() 
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- user table
CREATE TABLE "user" (
  "id"         UUID          NOT NULL,
  "email"      VARCHAR(255)  NOT NULL,
  "first_name" VARCHAR(255)  NOT NULL,
  "last_name"  VARCHAR(255)  NOT NULL,
  "google_id"  VARCHAR(255)  NOT NULL,
  "created_at" TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_user_email ON "user" ("email");
CREATE UNIQUE INDEX idx_user_google_id ON "user" ("google_id");

CREATE TRIGGER update_user_updated_at
    BEFORE UPDATE ON "user"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

END;
