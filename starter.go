package main

import (
	"bm.parts.server/config"
	"bm.parts.server/db"
	"bm.parts.server/log"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"strconv"
	"time"
)

type StRow struct {
	IdOC    int64
	IdProm  int64
	Price   float64
	InStock int8
}

type StList map[string]StRow

func main() {
	config.Start()
	log.Start()
	defer log.Stop()

	db.Start()
	//defer db.Stop()

	file, err := os.ReadFile(config.DirRoot + "/data/StartList.json")
	if err != nil {
		fmt.Println(err)
	}

	var data StList
	data = make(StList)
	_ = json.Unmarshal(file, &data)

	db.RebuildProductShort()

	count := 0
	now := time.Now().Unix()
	for idBM, row := range data {
		count++
		fmt.Println(strconv.Itoa(count), "idBM ", idBM, " row ", row)

		newRow := db.Short{
			ID:      primitive.NewObjectID(),
			IdBM:    idBM,
			IdOC:    row.IdOC,
			IdProm:  row.IdProm,
			Price:   row.Price,
			InStock: row.InStock,
			Updated: now,
		}

		result, err := db.ProductShort.InsertOne(context.TODO(), newRow)
		if err != nil {
			log.MyLog.Println(err)
		}

		log.MyLog.Println(result)
	}
}
