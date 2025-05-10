package object

type RoleCode string

const (
	RoleCodeAdmin         RoleCode = "SYSTEM_ADMIN"
	RoleCodeNormalUser    RoleCode = "GENERAL_USER"
	RoleCodeBusinessUser  RoleCode = "BUSINESS_USER"
	RoleCodeAccoutingUser RoleCode = "ACCOUNTING_USER"
)

func (rc RoleCode) String() string {
	return string(rc)
}
