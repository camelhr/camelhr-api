-- getUserByOrgSubdomainAPITokenQuery
-- $1: org_subdomain
-- $2: api_token
SELECT
    u.user_id,
    u.organization_id,
    u.email,
    u.password_hash,
    u.api_token,
    u.is_owner,
    is_email_verified,
    u.disabled_at,
    u.comment,
    u.created_at,
    u.updated_at,
    u.deleted_at
FROM
    users u
    JOIN organizations o ON u.organization_id = o.organization_id
WHERE
    o.subdomain = $1
    AND u.api_token = $2
    AND u.deleted_at IS NULL;
