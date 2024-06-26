-- +goose Up
-- +goose StatementBegin
CREATE TABLE organizations (
    organization_id SERIAL PRIMARY KEY,
    subdomain VARCHAR(30) NOT NULL UNIQUE CHECK (subdomain <> ''),
    name VARCHAR(60) NOT NULL UNIQUE CHECK (name <> ''),
    suspended_at TIMESTAMP WITHOUT TIME ZONE,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
    deleted_at TIMESTAMP WITHOUT TIME ZONE,
    comment VARCHAR(255)
);

-- create indexes
CREATE INDEX idx_organizations_deleted_at ON organizations(deleted_at);

-- create triggers to forbid truncate and delete operations on the organizations table
CREATE TRIGGER prevent_truncate_on_organizations
BEFORE TRUNCATE ON organizations
FOR EACH STATEMENT
EXECUTE FUNCTION operation_not_allowed();

CREATE TRIGGER prevent_hard_delete_on_organizations
BEFORE DELETE ON organizations
FOR EACH ROW
EXECUTE FUNCTION operation_not_allowed();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS organizations;
-- +goose StatementEnd
