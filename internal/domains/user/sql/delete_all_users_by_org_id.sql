-- deleteAllUsersByOrgIDQuery
-- $1 - organization_id
UPDATE
    users
SET
    deleted_at = NOW()
WHERE
    organization_id = $1
    AND deleted_at IS NULL
