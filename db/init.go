package db

import (
	"fsbm/conf"
)

var (
	FsbmSession *Handler
	migrations  map[string]func()
)

func Init() {
	FsbmSession = NewHandler()
	mysqlCfg := conf.GlobalConfig.Mysql
	FsbmSession.user = mysqlCfg.User
	FsbmSession.password = mysqlCfg.Password
	FsbmSession.ip = mysqlCfg.Ip
	FsbmSession.port = mysqlCfg.Port
	FsbmSession.dbName = mysqlCfg.DbName
	RunMigrations()
}

// 注册各表迁移函数
func RegisterMigration(table string, migration func()) {
	if migrations == nil {
		migrations = make(map[string]func())
	}
	migrations[table] = migration
}

func RunMigrations() {
	for _, fun := range migrations {
		fun()
	}
}
