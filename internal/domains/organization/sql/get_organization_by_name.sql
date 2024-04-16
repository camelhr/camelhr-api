-- getOrganizationByNameQuery
-- $1: name
SELECT
    organization_id,
    name,
    created_at,
    updated_at,
    deleted_at
FROM
    organization
WHERE
    name = $1;