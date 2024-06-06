package product

import (
	"bm.parts.server/config"
	"bm.parts.server/db"
	"bm.parts.server/dbok"
	"bm.parts.server/log"
	"context"
	"encoding/csv"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

type OKProduct struct {
	IdOk         int
	IdBM         string
	Model        string
	NameUk       string
	NameRu       string
	Price        float64
	Stock        int
	Manufacturer int
	Category     int
	Quantity     int
	Image        string
}

var Price map[string]float64
var OverPrice = 1.25
var Map map[string]db.Short
var CategoryMap map[string]int
var ManufacturerMap map[string]int

func PrepareProductShort() {
	now := time.Now().Unix()
	filter := bson.D{{}}
	update := bson.D{
		{"$set", bson.D{
			{"in_stock", 0},
			{"updated", now},
			{"new", 0},
		},
		},
	}
	_, err := db.ProductShort.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Prepare Product Short Done")
}

func UpdateFromOK(withDelete bool) {
	now := time.Now().Unix()
	filter := bson.D{{}}
	var update = bson.D{}
	if withDelete {
		update = bson.D{
			{"$set", bson.D{
				{"id_oc", 0},
				{"in_stock", 0},
				{"updated", now},
				{"new", 0},
			},
			},
		}
	} else {
		update = bson.D{
			{"$set", bson.D{
				{"updated", now},
				{"new", 0},
			},
			},
		}
	}
	_, err := db.ProductShort.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		fmt.Println(err)
	}

	client, _ := dbok.NewClient()
	defer client.Close()

	if withDelete {
		dbok.DisableAllProducts(client)
	}

	var productId int
	var sku string

	sql := "SELECT `product_id`, `sku` FROM `oc_product`"
	rows, err := client.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		err := rows.Scan(&productId, &sku)
		if err != nil {
			fmt.Println(err)
		}

		count++
		fmt.Println("From OK; Count", count, "Product Id", productId, "Sku", sku)

		idOC := productId
		idBM := sku

		now := time.Now().Unix()

		filter := bson.D{{"id_bm", idBM}}
		update := bson.D{
			{"$set", bson.D{
				{"id_oc", idOC},
				{"updated", now},
			},
			},
		}
		_, err = db.ProductShort.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func Update() {
	Price = make(map[string]float64)

	makeProductMap()
	makeCategoryMap()
	makeManufacturerMap()

	file, err := os.Open(config.DirRoot + "/data/BMPrice.csv")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	now := time.Now().Unix()
	count := 0
	for i, line := range data {
		count++
		fmt.Println("From Data; Count ", count)

		if i > 1 {
			brand := line[2]
			if brand == "MOTUL" {
				continue
			}

			idBM := line[0]
			price, err := strconv.ParseFloat(line[4], 64)
			if err != nil {
				fmt.Println(err)
			}
			price = price * OverPrice
			Price[idBM] = price

			oldProduct, ok := Map[idBM]
			if ok {
				if oldProduct.Price != price || oldProduct.InStock == 0 {
					filter := bson.D{{"_id", oldProduct.ID}}
					update := bson.D{
						{"$set", bson.D{
							{"price", price},
							{"in_stock", 1},
							{"updated", now},
						},
						},
					}
					_, err := db.ProductShort.UpdateOne(context.TODO(), filter, update)
					if err != nil {
						fmt.Println(err)
					}

					oldProduct.Price = price
					oldProduct.InStock = 1
					oldProduct.Updated = now
					Map[idBM] = oldProduct
				}
			} else {
				newRow := db.Short{
					ID:      primitive.NewObjectID(),
					IdBM:    idBM,
					IdOC:    0,
					IdProm:  0,
					Price:   price,
					InStock: 1,
					Updated: now,
					New:     1,
				}

				result, err := db.ProductShort.InsertOne(context.TODO(), newRow)
				if err != nil {
					log.MyLog.Println(err)
				}

				log.MyLog.Println(result)
			}
		}
	}
}

func makeProductMap() {
	Map = make(map[string]db.Short)

	cursor, err := db.ProductShort.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.MyLog.Println(err)
	}
	var res []db.Short
	if err = cursor.All(context.TODO(), &res); err != nil {
		log.MyLog.Println(err)
	}

	for i, result := range res {
		fmt.Println("Make Map; Count ", i)

		Map[result.IdBM] = result
	}
}

func makeCategoryMap() {
	CategoryMap = make(map[string]int)

	client, _ := dbok.NewClient()
	defer client.Close()

	var categoryId int
	var path string

	sql := "SELECT `category_id`, `path` FROM `oc_category_path_ua` ORDER BY `path`"
	rows, err := client.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		fmt.Println("From OK.CategoryPath; Count ", count)

		err := rows.Scan(&categoryId, &path)
		if err != nil {
			fmt.Println(err)
			continue
		}

		CategoryMap[path] = categoryId
	}
}

func makeManufacturerMap() {
	ManufacturerMap = make(map[string]int)

	client, _ := dbok.NewClient()
	defer client.Close()

	var brandId int
	var brandName string

	sql := "SELECT `manufacturer_id`, `name` FROM `oc_manufacturer` ORDER BY `name`"
	rows, err := client.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		fmt.Println("From OK.Manufacturer; Count ", count)

		err := rows.Scan(&brandId, &brandName)
		if err != nil {
			fmt.Println(err)
			continue
		}

		ManufacturerMap[brandName] = brandId
	}
}

func UpdateOK() {
	minDate := time.Now().Unix() - 12*60*60
	//minDate := time.Now().Unix() - 3*24*60*60

	threadsNumber := 12

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"id_oc", bson.D{{"$gt", 0}}}},
				bson.D{{"updated", bson.D{{"$gt", minDate}}}},
			}},
	}
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
	if total < 1000 {
		wg.Add(1)
		go updatePart(&wg, res)
	} else {
		limit := int(math.Ceil(float64(total / threadsNumber)))
		wg.Add(threadsNumber)
		fmt.Println("wg.Add 8")
		for i := 0; i < threadsNumber; i++ {
			min := i * limit
			max := (i + 1) * limit
			go updatePart(&wg, res[min:max])
		}
	}
	wg.Wait()
}

func updatePart(wg *sync.WaitGroup, rows []db.Short) {
	defer wg.Done()
	defer fmt.Println("wg.Done")

	client, _ := dbok.NewClient()
	defer client.Close()

	for i, row := range rows {
		dbok.ProductUpdate(client, row.IdOC, row.Price, row.InStock)
		fmt.Println("Update OK; Count ", i, "; Row.IdOC ", row.IdOC, "; Row.Price ", row.Price, " Row.InStock ", row.InStock)
	}
}
