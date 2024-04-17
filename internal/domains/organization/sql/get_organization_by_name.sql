-- getOrganizationByNameQuery
-- $1: name
SELECT
    organization_id,
    name,
    created_at,
    updated_at,
    deleted_at
FROM
    organizations
WHERE
    name = $1;