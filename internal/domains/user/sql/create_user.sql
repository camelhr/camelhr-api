-- createUserQuery
-- $1 - organization_id
-- $2 - email
-- $3 - password_hash
-- $4 - is_owner
INSERT INTO
    users(organization_id, email, password_hash, is_owner)
VALUES
    ($1, $2, $3, $4) RETURNING
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
    deleted_at;
