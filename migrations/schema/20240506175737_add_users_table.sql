-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    organization_id INTEGER NOT NULL,
    email VARCHAR(255) NOT NULL CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    password_hash TEXT,
    api_token text UNIQUE,
    is_owner BOOLEAN NOT NULL DEFAULT FALSE,
    is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    disabled_at TIMESTAMP WITHOUT TIME ZONE,
    comment VARCHAR(255),
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
    deleted_at TIMESTAMP WITHOUT TIME ZONE,
    FOREIGN KEY (organization_id) REFERENCES organizations(organization_id)
);

-- create unique constraint on organization_id and email
ALTER TABLE users ADD CONSTRAINT unique_organization_id_email UNIQUE (organization_id, email);
-- create partial unique index to ensure one owner per organization
CREATE UNIQUE INDEX idx_users_owner_per_organization ON users(organization_id) WHERE is_owner = TRUE;

-- create indexes
CREATE INDEX idx_users_organization_id ON users(organization_id);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- create triggers to forbid truncate and delete operations on the users table
CREATE TRIGGER prevent_truncate_on_users
BEFORE TRUNCATE ON users
FOR EACH STATEMENT
EXECUTE FUNCTION operation_not_allowed();

CREATE TRIGGER prevent_delete_on_users
BEFORE DELETE ON users
FOR EACH ROW
EXECUTE FUNCTION operation_not_allowed();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
