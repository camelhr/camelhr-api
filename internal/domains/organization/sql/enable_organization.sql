-- enableOrganizationQuery
-- $1: organization_id
-- $2: comment
UPDATE
    organizations
SET
    disabled_at = NULL,
    comment = $2
WHERE
    organization_id = $1
    AND disabled_at IS NOT NULL
    AND deleted_at IS NULL;
