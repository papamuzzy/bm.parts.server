package main

import (
	"bm.parts.server/apiBM"
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

	product.PrepareProductShort()
	apiBM.GetPrice()
	product.Update()
	product.UpdateOK()
	product.AddAll()
	product.UpdateFromOK(false)
}
