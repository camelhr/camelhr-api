-- updateOrganizationQuery
-- $1: organization_id
-- $2: name
UPDATE
    organizations
SET
    name = $2,
    updated_at = NOW()
WHERE
    organization_id = $1
    AND deleted_at IS NULL;
