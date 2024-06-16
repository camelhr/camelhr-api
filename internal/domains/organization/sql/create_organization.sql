-- createOrganizationQuery
-- $1: subdomain
-- $2: name
INSERT INTO
    organizations(subdomain, name)
VALUES
    ($1, $2) RETURNING
    organization_id,
    subdomain,
    name,
    suspended_at,
    created_at,
    updated_at,
    deleted_at,
    comment
