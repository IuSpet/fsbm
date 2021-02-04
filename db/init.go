package db

var fsbmSession *Handler

func Init() {
	fsbmSession = NewHandler()
	fsbmSession.user = "root"
	fsbmSession.password = "123456"
	fsbmSession.ip = "123456"
	fsbmSession.port = "3306"
	fsbmSession.dbName = "fsbm_test"
}
