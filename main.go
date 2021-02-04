package main

import (
	"fsbm/db"
	"fsbm/server"
)

func main() {
	db.Init()
	server.Run()
}
