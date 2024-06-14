package user

import _ "embed"

//go:embed sql/get_user_by_id.sql
var getUserByIDQuery string

//go:embed sql/get_user_by_api_token.sql
var getUserByAPITokenQuery string

//go:embed sql/get_user_by_org_subdomain_api_token.sql
var getUserByOrgSubdomainAPITokenQuery string

//go:embed sql/get_user_by_org_id_email.sql
var getUserByOrgIDEmailQuery string

//go:embed sql/get_user_by_org_subdomain_email.sql
var getUserByOrgSubdomainEmailQuery string

//go:embed sql/create_user.sql
var createUserQuery string

//go:embed sql/reset_password.sql
var resetPasswordQuery string

//go:embed sql/delete_user.sql
var deleteUserQuery string

//go:embed sql/delete_all_users_by_org_id.sql
var deleteAllUsersByOrgIDQuery string

//go:embed sql/disable_user.sql
var disableUserQuery string

//go:embed sql/enable_user.sql
var enableUserQuery string

//go:embed sql/genereate_api_token.sql
var generateAPITokenQuery string

//go:embed sql/reset_api_token.sql
var resetAPITokenQuery string

//go:embed sql/set_email_verified.sql
var setEmailVerifiedQuery string
