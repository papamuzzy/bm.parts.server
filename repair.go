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

	//product.UpdateFromOK()
	//product.GetPrice()
	//product.Update()
	//product.UpdateOK()
	product.RepairFromOK()
	product.AddAll()
	product.UpdateFromOK()
}
