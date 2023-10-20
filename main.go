package main

import (
	"opa-test/config"
	"opa-test/routers"
)

func main() {
	db := config.ConnectionDatabase()
	routers.Init(db)
}
