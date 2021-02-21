package server

import "fsbm/util/auth"

// 接口权限配置
var AllPathPermission = map[string][]auth.SubjectPermission{
	"/admin/user_list": {
		{
			PermissionType: "api",
			PermissionName: "user_list",
		},
	},
	"admin/authority/modify": {
		{
			PermissionType: "api",
			PermissionName: "authority_modify",
		},
	},
	"admin/user_detail": {
		{
			PermissionType: "api",
			PermissionName: "user_detail",
		},
	},
}
