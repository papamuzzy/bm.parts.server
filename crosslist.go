package main

import (
	"bm.parts.server/config"
	"bm.parts.server/db"
	"bm.parts.server/log"
)

func main() {
	config.Start()
	log.Start()
	defer log.Stop()

	db.Start()
	//defer db.Stop()

}
