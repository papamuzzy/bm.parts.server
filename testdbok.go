package main

import (
	"bm.parts.server/config"
	"bm.parts.server/db"
	"bm.parts.server/log"
	"bm.parts.server/product"
)

func main() {
	config.Start()
	log.Start()
	defer log.Stop()

	db.Start()
	//defer db.Stop()

	product.UpdateProductCat("BFFFF4653B35EEA047CA7F5DBBE3AF27")
}
