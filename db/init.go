package db

import "fsbm/conf"

var fsbmSession *Handler

//func init() {
//	Init()
//}

func Init() {
	fsbmSession = NewHandler()
	mysqlCfg := conf.GlobalConfig.Mysql
	fsbmSession.user = mysqlCfg.User
	fsbmSession.password = mysqlCfg.Password
	fsbmSession.ip = mysqlCfg.Ip
	fsbmSession.port = mysqlCfg.Port
	fsbmSession.dbName = mysqlCfg.DbName
}
