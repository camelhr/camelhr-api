-- getUserByOrgIDEmailQuery
-- $1: org_id
-- $2: email
SELECT
    user_id,
    organization_id,
    email,
    password_hash,
    api_token,
    is_owner,
    disabled_at,
    comment,
    created_at,
    updated_at,
    deleted_at
FROM
    users
WHERE
    organization_id = $1
    AND email = $2
    AND deleted_at IS NULL;
