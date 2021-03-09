package admin

import (
	"fsbm/db"
	"time"
)

func getUserList(name, email, phone string, gender, age int8, begin, end time.Time, page, pageSize int64) (res []db.UserAccountInfo, err error) {
	conn, err := db.FsbmSession.GetConnection()
	if err != nil {
		return
	}
	conn = conn.Where("created_at > ? and created_at < ?", begin, end)
	if name != "" {
		conn = conn.Where("name like ?", "%"+name+"%")
	}
	if email != "" {
		conn = conn.Where("email = ?", email)
	}
	if phone != "" {
		conn = conn.Where("phone = ?", phone)
	}
	if gender != -1 {
		conn = conn.Where("gender = ?", gender)
	}
	if age != -1 {
		conn = conn.Where("age = ?", age)
	}
	if page != -1 {
		limit := int(pageSize)
		offset := int(pageSize * (page - 1))
		conn = conn.Limit(limit).Offset(offset)
	}
	err = conn.Debug().Find(&res).Error
	return
}
