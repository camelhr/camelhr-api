-- +goose Up
-- +goose StatementBegin
CREATE TABLE organizations (
    organization_id SERIAL PRIMARY KEY,
    name VARCHAR(60) NOT NULL UNIQUE,
    suspended_at TIMESTAMP WITHOUT TIME ZONE,
    blacklisted_at TIMESTAMP WITHOUT TIME ZONE,
    comment VARCHAR(255),
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE organizations;
-- +goose StatementEnd
