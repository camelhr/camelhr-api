-- disableOrganizationQuery
-- $1: organization_id
-- $2: comment
UPDATE
    organizations
SET
    disabled_at = NOW(),
    comment = $2
WHERE
    organization_id = $1
    AND disabled_at IS NULL
    AND deleted_at IS NULL;
