package db

var fsbmSession *Handler

func init() {
	Init()
}

func Init() {
	fsbmSession = NewHandler()
	fsbmSession.user = "root"
	fsbmSession.password = "123456"
	fsbmSession.ip = "127.0.0.1"
	fsbmSession.port = "3306"
	fsbmSession.dbName = "fsbm_test"
}
