package cronjob

import (
	"context"
	"fsbm/db"
	"fsbm/util/logs"
	"time"
)

// 检查所有权限
func authCheckTask(ctx context.Context) error {
	relationList, err := db.GetAllActiveRelation()
	if err != nil {
		logs.CtxError(ctx, "get all active relation error. err: %+v", err)
		return err
	}
	now := time.Now()
	for idx := range relationList {
		if now.After(relationList[idx].EndTime) {
			relationList[idx].Status = db.AuthUserRoleStatus_Expired
			logs.CtxInfo(ctx, "user[%d]'s role[%d] expired", relationList[idx].UserID, relationList[idx].RoleID)
		}
	}
	err = db.SaveAuthUserRoleRows(relationList)
	if err != nil {
		logs.CtxError(ctx, "save user role relation error. err: %+v", err)
		return err
	}
	return nil
}
