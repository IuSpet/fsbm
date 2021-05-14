package main

import (
	"fsbm/conf"
	"fsbm/cronjob"
	"fsbm/db"
	"fsbm/server"
)

func main() {
	conf.Init()
	db.Init()
	cronjob.RunCronJob()
	server.Run()
}
