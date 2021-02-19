package auth

import (
	"context"
	"fsbm/db"
)

type UserRoleSubject struct {
	Email          string
	Role           []db.AuthRole
	PermissionList map[string]bool // [type+name] -> true
}

// 根据用户唯一标示(邮箱)获取用户权限信息
func NewUserRoleSubject(email string) (*UserRoleSubject, error) {
	user, err := db.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	roleList, err := db.GetRoleById(user.ID)
	if err != nil {
		return nil, err
	}
	var roleIDList []int64
	for _, role := range roleList {
		roleIDList = append(roleIDList, role.ID)
	}
	permissionList, err := db.GetPermissionByRoleID(roleIDList)
	if err != nil {
		return nil, err
	}
	permissionMap := make(map[string]bool)
	for _, permission := range permissionList {
		permissionMap[permission.String()] = true
	}
	return &UserRoleSubject{
		Email:          email,
		Role:           roleList,
		PermissionList: permissionMap,
	}, nil
}

// 判断对页面/接口是否有权限
func (u *UserRoleSubject) HasPermission(ctx context.Context, permissionList []SubjectPermission) bool {
	hasPermission := true
	for _, permission := range permissionList {
		if u.PermissionList[permission.String()] {
			continue
		}
		hasPermission = false
		break
	}
	return hasPermission
}
