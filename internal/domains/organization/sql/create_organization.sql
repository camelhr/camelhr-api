-- createOrganizationQuery
-- $1: name
-- $2: description
INSERT INTO
    organizations(name, description)
VALUES
($1, $2) RETURNING organization_id;