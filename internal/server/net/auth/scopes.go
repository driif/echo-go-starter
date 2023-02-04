package auth

type Scope string

const (
	AuthScopeApp        Scope = "app"
	AuthScopeAdmin      Scope = "admin"
	AuthScopeSuperAdmin Scope = "superadmin"
)

func (s Scope) String() string {
	return string(s)
}
