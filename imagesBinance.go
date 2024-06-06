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

	//product.UpdateFromOK2()
	//product.GetPrice()
	//product.Update3()
	//product.UpdateOK()
	product.AddImg()
	//product.AddAllToMongo()
	//product.UpdateFromOK()
}
