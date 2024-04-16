-- createOrganizationQuery
-- $1: name
-- $2: description
INSERT INTO
    organization(name, description)
VALUES
($1, $2) RETURNING organization_id;