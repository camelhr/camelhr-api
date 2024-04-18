-- deleteOrganizationQuery
-- $1: organization_id
UPDATE
    organizations
SET
    deleted_at = NOW()
WHERE
    organization_id = $1 
    AND deleted_at IS NULL;
