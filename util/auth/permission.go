package auth

type SubjectPermission struct {
	PermissionType string
	PermissionName string
}

func (s SubjectPermission) String() string {
	return s.PermissionType + ":" + s.PermissionName
}
