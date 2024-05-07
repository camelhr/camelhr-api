-- updateOrganizationQuery
-- $1: organization_id
-- $2: subdomain
-- $3: name
UPDATE
    organizations
SET
    subdomain = $2,
    name = $3,
    updated_at = NOW()
WHERE
    organization_id = $1
    AND deleted_at IS NULL;
