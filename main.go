package main

import (
	"fsbm/conf"
	"fsbm/db"
	"fsbm/server"
)

func main() {
	conf.Init()
	db.Init()
	server.Run()
}
