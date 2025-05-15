BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

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
  "id"                    UUID         NOT NULL,
  "email"                 VARCHAR(255) NOT NULL,
  "first_name"            VARCHAR(255) NOT NULL,
  "last_name"             VARCHAR(255) NOT NULL,
  "refresh_token_hash"    VARCHAR(255) NOT NULL,
  "google_id"             VARCHAR(255) NOT NULL,
  "scheduled_deletion_at" TIMESTAMPTZ  DEFAULT NULL,
  "created_at"            TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"            TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
  "id"             UUID         NOT NULL,
  "user_id"        UUID         NOT NULL,
  "key_hash"       VARCHAR(255) NOT NULL,
  "key_ciphertext" BYTEA NOT NULL,
  "key_nonce"      BYTEA NOT NULL,
  "created_at"     TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"     TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_license_key_hash ON "license" ("key_hash");
CREATE UNIQUE INDEX idx_license_key_ciphertext_nonce ON "license" ("key_ciphertext", "key_nonce");

CREATE TRIGGER update_license_updated_at
    BEFORE UPDATE ON "license"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- plan table
CREATE TABLE "plan" (
  "id"              UUID          NOT NULL,
  "name"            VARCHAR(255)  NOT NULL,
  "price"           NUMERIC(10,2) NOT NULL,
  "stripe_price_id" VARCHAR(255)  NOT NULL,
  "created_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"      TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_plan_stripe_price_id ON "plan" ("stripe_price_id");

CREATE TRIGGER update_plan_updated_at
    BEFORE UPDATE ON "plan"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

INSERT INTO "plan" ("id", "name", "price", "stripe_price_id") VALUES
  (gen_random_uuid(), 'Team', 10, 'team'),
  (gen_random_uuid(), 'Business', 24, 'business');

-- subscription table
CREATE TABLE "subscription_status" (
  "code" INTEGER      NOT NULL,
  "name" VARCHAR(255) NOT NULL,
  PRIMARY KEY ("code")
);

CREATE UNIQUE INDEX idx_subscription_status_code ON "subscription_status" ("code");
CREATE UNIQUE INDEX idx_subscription_status_name ON "subscription_status" ("name");

INSERT INTO "subscription_status" ("code", "name") VALUES
  (0, 'unknown'),
  (1, 'trial'),
  (2, 'active'),
  (3, 'canceled'),
  (4, 'past_due');

CREATE TABLE "subscription" (
  "id"                     UUID         NOT NULL,
  "user_id"                UUID         NOT NULL,
  "plan_id"                UUID,
  "stripe_customer_id"     VARCHAR(255) NOT NULL,
  "stripe_subscription_id" VARCHAR(255) NOT NULL,
  "status"                 INTEGER      NOT NULL,
  "trial_start"            TIMESTAMPTZ  NOT NULL,
  "trial_end"              TIMESTAMPTZ  NOT NULL,
  "seat_count"             INTEGER      NOT NULL DEFAULT 1,
  "created_at"             TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at"             TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("plan_id") REFERENCES "plan" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("status") REFERENCES "subscription_status" ("code") ON DELETE RESTRICT,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX idx_subscription_stripe_customer_id ON "subscription" ("stripe_customer_id");
CREATE UNIQUE INDEX idx_subscription_stripe_subscription_id ON "subscription" ("stripe_subscription_id");

CREATE TRIGGER update_subscription_updated_at
    BEFORE UPDATE ON "subscription"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
END;