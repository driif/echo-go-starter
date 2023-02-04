package status

type UserStatus int

const (
	Active UserStatus = iota
	Inactive
	NotDisturb
	Offline
)

func (s UserStatus) String() string {
	return [...]string{"active", "inactive", "not_disturb", "offline"}[s]
}
