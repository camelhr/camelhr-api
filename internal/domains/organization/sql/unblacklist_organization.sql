-- unblacklistOrganizationQuery
-- $1: organization_id
-- $2: comment
UPDATE
    organization
SET
    blacklisted_at = NULL,
    comment = $2
WHERE
    organization_id = $1
    AND blacklisted_at IS NOT NULL
    AND deleted_at IS NULL;
