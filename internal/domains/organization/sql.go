package organization

import _ "embed"

//go:embed sql/get_organization_by_id.sql
var getOrganizationByIDQuery string

//go:embed sql/get_organization_by_subdomain.sql
var getOrganizationBySubdomainQuery string

//go:embed sql/get_organization_by_name.sql
var getOrganizationByNameQuery string

//go:embed sql/create_organization.sql
var createOrganizationQuery string

//go:embed sql/update_organization.sql
var updateOrganizationQuery string

//go:embed sql/delete_organization.sql
var deleteOrganizationQuery string

//go:embed sql/suspend_organization.sql
var suspendOrganizationQuery string

//go:embed sql/unsuspend_organization.sql
var unsuspendOrganizationQuery string

//go:embed sql/disable_organization.sql
var disableOrganizationQuery string

//go:embed sql/enable_organization.sql
var enableOrganizationQuery string
