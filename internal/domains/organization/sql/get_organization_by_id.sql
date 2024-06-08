-- getOrganizationByIDQuery
-- $1: organization_id
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
    organization_id = $1
    AND deleted_at IS NULL;
