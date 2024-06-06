package product

import (
	"bm.parts.server/apiBM"
	"bm.parts.server/apiOK"
	"bm.parts.server/db"
	"bm.parts.server/log"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"math"
	"sync"
)

func PhotoAll() {
	threadsNumber := 12

	filter := bson.D{{"id_oc", bson.D{{"$ne", 0}}}}
	cursor, err := db.ProductShort.Find(context.TODO(), filter)
	if err != nil {
		log.MyLog.Println(err)
	}
	var res []db.Short
	if err = cursor.All(context.TODO(), &res); err != nil {
		log.MyLog.Println(err)
	}

	total := len(res)
	var wg sync.WaitGroup
	if total < 100 {
		wg.Add(1)
		go photoPart(&wg, res)
	} else {
		limit := int(math.Ceil(float64(total / threadsNumber)))
		wg.Add(threadsNumber)
		fmt.Println("wg.Add 8")
		for i := 0; i < threadsNumber; i++ {
			min := i * limit
			max := (i + 1) * limit
			go photoPart(&wg, res[min:max])
		}
	}
	wg.Wait()
}

func photoPart(wg *sync.WaitGroup, rows []db.Short) {
	defer wg.Done()
	defer fmt.Println("wg.Done")

	for i, row := range rows {
		idBM := row.IdBM

		_ = AddPhoto(idBM)

		fmt.Println("Add OK; Count ", i, " Row ", row)
	}
}

func AddPhoto(idBM string) bool {
	product, ok := apiBM.GetProduct(idBM)
	if ok {
		resp := apiOK.Post("https://dev.dmdshop.com.ua/index.php?route=tool/add/photo", product)
		fmt.Println(resp.Status)

		return true
	}

	return false
}
