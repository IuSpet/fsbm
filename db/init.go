package db

import (
	"fsbm/conf"
)

var (
	fsbmSession *Handler
	migrations  map[string]func()
)

func Init() {
	fsbmSession = NewHandler()
	mysqlCfg := conf.GlobalConfig.Mysql
	fsbmSession.user = mysqlCfg.User
	fsbmSession.password = mysqlCfg.Password
	fsbmSession.ip = mysqlCfg.Ip
	fsbmSession.port = mysqlCfg.Port
	fsbmSession.dbName = mysqlCfg.DbName
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
