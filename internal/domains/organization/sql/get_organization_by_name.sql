-- getOrganizationByNameQuery
-- $1: name
SELECT
    organization_id,
    subdomain,
    name,
    suspended_at,
    disabled_at,
    comment,
    created_at,
    updated_at,
    deleted_at
FROM
    organizations
WHERE
    name = $1
    AND deleted_at IS NULL;
