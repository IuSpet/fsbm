package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Handler struct {
	db       *gorm.DB
	config   *gorm.Config
	user     string
	password string
	ip       string
	port     string
	dbName   string
}

func NewHandler() *Handler {
	return &Handler{config: &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}}
}

func (h *Handler) GetConnection() (*gorm.DB, error) {
	var err error
	if h.db != nil {
		return h.db, nil
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		h.user, h.password, h.ip, h.port, h.dbName)
	h.db, err = gorm.Open(mysql.Open(dsn), h.config)
	if err != nil {
		return nil, err
	}
	return h.db, nil
}
