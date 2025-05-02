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
  "id"                 UUID          NOT NULL,
  "email"              VARCHAR(255)  NOT NULL,
  "first_name"         VARCHAR(255)  NOT NULL,
  "last_name"          VARCHAR(255)  NOT NULL,
  "refresh_token_hash" VARCHAR(255)  NOT NULL,
  "google_id"          VARCHAR(255)  NOT NULL,
  "created_at"         TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"         TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_user_email ON "user" ("email");
CREATE UNIQUE INDEX idx_user_refresh_token_hash ON "user" ("refresh_token_hash");

CREATE TRIGGER update_user_updated_at
    BEFORE UPDATE ON "user"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- license table
CREATE TABLE "license" (
  "id"         UUID         NOT NULL,
  "user_id"    UUID         NOT NULL,
  "key"        VARCHAR(255) NOT NULL,
  "created_at" TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_license_key ON "license" ("key");

CREATE TRIGGER update_license_updated_at
    BEFORE UPDATE ON "license"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

END;