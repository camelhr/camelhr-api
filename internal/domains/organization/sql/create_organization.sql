-- createOrganizationQuery
-- $1: name
INSERT INTO
    organizations(name)
VALUES
    ($1) RETURNING organization_id;
