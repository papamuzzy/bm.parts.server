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
	product.PhotoAll()
	//product.AddPhoto("A467006631399E3D4D11ADF44B13ECC2")
}
