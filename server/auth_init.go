package server

import "fsbm/util/auth"

var AllPathPermission = map[string][]auth.SubjectPermission{
	"/ping":[]auth.SubjectPermission{
		{
			PermissionType: "none",
			PermissionName: "ping",
		},
	},
}

