package task

import (
	"fsbm/db"
	"time"
)

func SetUserOperationByEmail(email, op string) (err error) {
	var user *db.UserAccountInfo
	for i := 0; i < 3; i++ {
		user, err = db.GetUserByEmail(email)
		if err == nil {
			break
		}
	}
	if err != nil {
		return err
	}
	err = SetUserOperationById(user.ID, op)
	return
}

func SetUserOperationById(id int64, op string) (err error) {
	now := time.Now()
	row := db.UserOperationLog{
		UserId:     id,
		Operation:  op,
		OperatedAt: now.Unix(),
	}
	for i := 0; i < 3; i++ {
		err = db.SaveUserOperationLog(&row)
		if err == nil {
			break
		}
	}
	return
}
