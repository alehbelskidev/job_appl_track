CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
  id       UUID        NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
  email    TEXT        NOT NULL UNIQUE,
  password TEXT        NOT NULL
);

CREATE TABLE IF NOT EXISTS job_applications (
  id           UUID        NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
  company      TEXT        NOT NULL,
  role         TEXT        NOT NULL,
  date_applied TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  date_updated TIMESTAMPTZ,
  status       INTEGER     NOT NULL DEFAULT 0,
  owner_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
  id           UUID        NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
  token TEXT   NOT NULL UNIQUE
);
