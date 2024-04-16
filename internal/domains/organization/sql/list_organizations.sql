-- listOrganizationsQuery

SELECT
    o.organization_id,
    o.name,
    o.created_at,
    o.updated_at,
    o.deleted_at
FROM
    organization o
WHERE
    o.deleted_at IS NULL
ORDER BY
    o.name;
