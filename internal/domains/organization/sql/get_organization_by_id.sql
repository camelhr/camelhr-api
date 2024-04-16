-- getOrganizationByIDQuery
-- $1: organization_id
SELECT
    organization_id,
    name,
    created_at,
    updated_at,
    deleted_at
FROM
    organization
WHERE
    organization_id = $1;