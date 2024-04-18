-- unsuspendOrganizationQuery
-- $1: organization_id
UPDATE
    organization
SET
    suspended_at = NULL
WHERE
    organization_id = $1
    AND deleted_at IS NULL;
