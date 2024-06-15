-- deleteOrganizationQuery
-- $1: organization_id
-- $2: comment
UPDATE
    organizations
SET
    deleted_at = NOW(),
    comment = $2
WHERE
    organization_id = $1 
    AND deleted_at IS NULL;
