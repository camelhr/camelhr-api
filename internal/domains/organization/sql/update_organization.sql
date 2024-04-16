-- updateOrganizationQuery
-- $1: organization_id
-- $2: name
UPDATE
    organization
SET
    name = $2,
    updated_at = NOW()
WHERE
    organization_id = $1;