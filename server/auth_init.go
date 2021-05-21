package server

import "fsbm/util/auth"

// 接口权限配置
var AllPathPermission = map[string][]auth.SubjectPermission{
	"/admin/user_list": {
		{
			PermissionType: "read",
			PermissionName: "user_list",
		},
	},
	"/admin/authority/modify": {
		{
			PermissionType: "write",
			PermissionName: "authority_modify",
		},
	},
	"/admin/user_detail": {
		{
			PermissionType: "api",
			PermissionName: "user_detail",
		},
	},
}
