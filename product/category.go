package product

import (
	"bm.parts.server/apiBM"
	"bm.parts.server/db"
	"bm.parts.server/dbok"
	"bm.parts.server/log"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"math"
	"strings"
	"sync"
)

func ProductsToCategory() {
	makeProductMap()

	client, _ := dbok.NewClient()
	defer client.Close()

	sql := "DELETE FROM `oc_product_to_category` WHERE ?"
	_, err := client.Exec(sql, 1)
	if err != nil {
		fmt.Println(err)
	}

	count := 0
	for idBM, line := range Map {
		count++
		fmt.Println("From Data; Count ", count)
		idOK := line.IdOC

		product, ok := apiBM.GetProduct(idBM)
		if ok {
			dataUK := product["UK"].(map[string]interface{})
			pathUK := dataUK["nodes"].(string)

			dataRU := product["RU"].(map[string]interface{})
			pathRU := dataRU["nodes"].(string)

			categoryId := GetCategory(client, pathUK, pathRU)

			fmt.Println("Category Id ", categoryId)

			sql = "INSERT INTO `oc_product_to_category` SET `product_id` = ?, `category_id` = ?, `main_category` = 1"
			_, err := client.Exec(sql, idOK, categoryId)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func CategoryAll() {
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
		go categoryPart(&wg, res)
	} else {
		limit := int(math.Ceil(float64(total / threadsNumber)))
		wg.Add(threadsNumber)
		fmt.Println("wg.Add 8")
		for i := 0; i < threadsNumber; i++ {
			min := i * limit
			max := (i + 1) * limit
			go categoryPart(&wg, res[min:max])
		}
	}
	wg.Wait()
}

func categoryPart(wg *sync.WaitGroup, rows []db.Short) {
	defer wg.Done()
	defer fmt.Println("wg.Done")

	for i, row := range rows {
		idBM := row.IdBM
		idOC := row.IdOC

		_ = AddCategory(idBM, int(idOC))

		fmt.Println("Add Category OC; Count ", i, " Row ", row)
	}
}

type dbOKPath struct {
	CategoryId int
	PathId     int
	Level      int
}

func AddCategory(idBM string, idOC int) bool {
	product, ok := apiBM.GetProduct(idBM)
	if ok {
		dataUK := product["UK"].(map[string]interface{})
		pathUK := strings.ToLower(dataUK["nodes"].(string))

		client, _ := dbok.NewClient()
		defer client.Close()

		sql := "SELECT `c`.`category_id` FROM `oc_category` AS `c` LEFT JOIN `oc_category_description` AS `cd` ON `cd`.`category_id` = `c`.`category_id` AND `cd`.`language_id` = 3 WHERE `cd`.`path` = ?"
		var catId int
		err := client.QueryRow(sql, pathUK).Scan(&catId)
		if err != nil || catId == 0 {
			fmt.Println(err, catId)
			return false
		}

		fmt.Println("Cat Id ", catId)

		sql = "INSERT INTO `oc_product_to_category` SET `product_id` = ?, `category_id` = ?, `main_category` = 1"
		_, err = client.Exec(sql, idOC, catId)

		return true
	}

	return false
}
