-- getOrganizationByNameQuery
-- $1: name
SELECT
    organization_id,
    subdomain,
    name,
    suspended_at,
    created_at,
    updated_at,
    deleted_at,
    comment
FROM
    organizations
WHERE
    name = $1
    AND deleted_at IS NULL;
