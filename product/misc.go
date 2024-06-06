package product

import (
	"bm.parts.server/config"
	"bm.parts.server/db"
	"bm.parts.server/dbok"
	"bm.parts.server/log"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var MapOK map[string]int

var ProductsOK map[string]OKProduct

func RepairFromOK() {
	makeCategoryMap()
	makeManufacturerMap()

	now := time.Now().Unix()
	filter := bson.D{{}}
	update := bson.D{
		{"$set", bson.D{
			{"updated", now},
			{"new", 1},
		},
		},
	}
	_, err := db.ProductShort.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		fmt.Println(err)
	}

	client, _ := dbok.NewClient()
	defer client.Close()

	var product_id int
	var sku string

	sql := "SELECT `product_id`, `sku` FROM `oc_product`"
	rows, err := client.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		fmt.Println("From OK; Count ", count)

		err := rows.Scan(&product_id, &sku)
		if err != nil {
			fmt.Println(err)
		}

		idBM := sku

		now := time.Now().Unix()

		filter := bson.D{{"id_bm", idBM}}
		update := bson.D{
			{"$set", bson.D{
				{"new", 0},
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

func UpdateFromOK2() {
	filter := bson.D{{}}
	update := bson.D{
		{"$set", bson.D{
			{"id_oc", 0},
		},
		},
	}
	_, err := db.ProductShort.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		fmt.Println(err)
	}

	client, _ := dbok.NewClient()
	defer client.Close()

	var product_id int
	var sku string

	sql := "SELECT `product_id`, `sku` FROM `oc_product`"
	rows, err := client.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		fmt.Println("From OK; Count ", count)

		err := rows.Scan(&product_id, &sku)
		if err != nil {
			fmt.Println(err)
		}

		idOC := product_id
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

func Update2() {
	Price = make(map[string]float64)

	makeProductMap()

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
			idBM := line[0]
			price, err := strconv.ParseFloat(line[4], 64)
			if err != nil {
				fmt.Println(err)
			}
			price = price * OverPrice
			Price[idBM] = price

			oldProduct, ok := Map[idBM]
			if ok {
				if empty := UpdateProductCat(idBM); empty == 1 {
					filter := bson.D{{"_id", oldProduct.ID}}
					update := bson.D{
						{"$set", bson.D{
							{"new", 1},
						},
						},
					}
					_, err := db.ProductShort.UpdateOne(context.TODO(), filter, update)
					if err != nil {
						fmt.Println(err)
					}
				} else if oldProduct.Price != price || oldProduct.InStock == 0 {
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

func Update3() {
	Price = make(map[string]float64)

	makeProductMap()

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

					filter = bson.D{{"id_bm", idBM}}
					update = bson.D{
						{"$set", bson.D{
							{"uk.price", price},
						},
						},
					}
					_, err = db.ProductFull.UpdateOne(context.TODO(), filter, update)
					if err != nil {
						fmt.Println(err)
					}

					oldProduct.Price = price
					oldProduct.InStock = 1
					oldProduct.Updated = now
					Map[idBM] = oldProduct
				}
			}
		}
	}
}

func UpdateCatAll() {
	makeCategoryMap()
	makeManufacturerMap()

	total, _ := db.ProductFull.EstimatedDocumentCount(context.TODO())

	cursor, err := db.ProductFull.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.MyLog.Println(err)
		return
	}

	count := 0
	for cursor.Next(context.TODO()) {
		count++
		fmt.Println("Update category & brand ", count, " From ", total)
		var res db.Full
		if err := cursor.Decode(&res); err != nil {
			fmt.Println(err)
			continue
		}

		UpdateProdCat(res)
	}
}

func UpdateCatAllOK() {
	makeProductMapOK()

	total, _ := db.ProductFull.EstimatedDocumentCount(context.TODO())

	cursor, err := db.ProductFull.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.MyLog.Println(err)
		return
	}

	count := 0
	for cursor.Next(context.TODO()) {
		count++
		fmt.Println("Update in OK category & brand ", count, " From ", total)
		var res db.Full
		if err := cursor.Decode(&res); err != nil {
			fmt.Println(err)
			continue
		}

		productId, ok := MapOK[res.IdBM]
		if ok {
			UpdateProdCatOK(res, productId)
		}
	}
}

func CheckProductsAllOK() {
	makeProductsOK()

	client, _ := dbok.NewClient()
	defer client.Close()

	total, _ := db.ProductFull.EstimatedDocumentCount(context.TODO())

	cursor, err := db.ProductFull.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.MyLog.Println(err)
		return
	}

	countNotChecked := 0
	countNotCheckedImg := 0
	countNotCheckedNameUk := 0
	countNotCheckedNameRu := 0
	countNotCheckedModel := 0
	countNotCheckedPrice := 0
	countNotCheckedBrand := 0
	countNotCheckedCategory := 0
	countNotCheckedStock := 0
	countNotCheckedQuantity := 0

	count := 0
	for cursor.Next(context.TODO()) {
		count++
		fmt.Println("Check product in OK ", count, " From ", total)
		var res db.Full
		if err := cursor.Decode(&res); err != nil {
			fmt.Println(err)
			continue
		}

		idBM := res.IdBM
		dataUk := res.UK
		nameUk := dataUk["name"]
		dataRu := res.RU
		nameRu := dataRu["name"]
		modelBM := dataUk["article"].(string)
		brandId := res.ManufacturerId
		catId := res.CategoryId
		price, _ := strconv.ParseFloat(dataUk["price"].(string), 64)
		price = price * OverPrice
		imageArr := strings.Split(dataUk["default_image"].(string), "\\")
		image := imageArr[len(imageArr)-1]

		productOK, ok := ProductsOK[idBM]
		if ok {
			productId := productOK.IdOk

			imageArr = strings.Split(productOK.Image, "/")
			imageOK := imageArr[len(imageArr)-1]

			checked := true

			if nameUk != productOK.NameUk {
				countNotCheckedNameUk++
				checked = false
				fmt.Println("BM Uk ", nameUk)
				fmt.Println("OK Uk ", productOK.NameUk)
				fmt.Println("Name Uk Update")
				sql := "UPDATE `oc_product_description` SET `name` = ? WHERE product_id = ? AND `language_id` = 1"
				_, err = client.Exec(sql, nameUk, productId)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Name Uk Updated!")
				}
			}
			if nameRu != productOK.NameRu {
				countNotCheckedNameRu++
				checked = false
				fmt.Println("BM Ru ", nameRu)
				fmt.Println("OK Ru ", productOK.NameRu)
				fmt.Println("Name Ru Update")
				sql := "UPDATE `oc_product_description` SET `name` = ? WHERE product_id = ? AND `language_id` = 3"
				_, err = client.Exec(sql, nameRu, productId)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Name Ru Updated!")
				}
			}
			if checked && modelBM != productOK.Model {
				countNotCheckedModel++
				checked = false
				fmt.Println("Model BM ", modelBM)
				fmt.Println("Model OK ", productOK.Model)
				fmt.Println("Model Update")
				sql := "UPDATE `oc_product` SET `model` = ? WHERE `product_id` = ?"
				_, err = client.Exec(sql, modelBM, productId)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Model Updated!")
				}
			}
			if checked && math.Abs(price-productOK.Price) > 0.009 {
				countNotCheckedPrice++
				checked = false
				fmt.Println("Price BM ", price)
				fmt.Println("Price OK ", productOK.Price)
				fmt.Println("Price Update")
				sql := "UPDATE `oc_product` SET `price` = ? WHERE `product_id` = ?"
				_, err = client.Exec(sql, price, productId)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Price Updated!")
				}
			}
			if checked && brandId != productOK.Manufacturer {
				countNotCheckedBrand++
				checked = false
				fmt.Println("Brand BM ", brandId)
				fmt.Println("Brand OK ", productOK.Manufacturer)
				fmt.Println("Brand Update")
				sql := "UPDATE `oc_product` SET `manufacturer_id` = ? WHERE `product_id` = ?"
				_, err = client.Exec(sql, brandId, productId)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Brand Updated!")
				}
			}
			if checked && catId != productOK.Category {
				countNotCheckedCategory++
				checked = false
				fmt.Println("Category BM ", catId)
				fmt.Println("Category OK ", productOK.Category)
				fmt.Println("Category Update")
				sql := "UPDATE `oc_product_to_category` SET `category_id` = ? WHERE `product_id` = ?"
				_, err = client.Exec(sql, catId, productId)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Category Updated!")
				}
			}
			if checked && image != imageOK {
				countNotCheckedImg++
				checked = false
			}
			if checked && productOK.Stock != 7 {
				countNotCheckedStock++
				checked = false
				fmt.Println("Stock OK ", productOK.Stock)
				fmt.Println("Stock Update")
				sql := "UPDATE `oc_product` SET `stock_status_id` = 7 WHERE `product_id` = ?"
				_, err = client.Exec(sql, productId)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Stock Updated!")
				}
			}
			if checked && productOK.Quantity == 0 {
				countNotCheckedQuantity++
				checked = false
				fmt.Println("Quantity OK ", productOK.Quantity)
				fmt.Println("Quantity Update")
				sql := "UPDATE `oc_product` SET `quantity` = 100 WHERE `product_id` = ?"
				_, err = client.Exec(sql, productId)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Quantity Updated!")
				}
			}

			if !checked {
				countNotChecked++

				/*fmt.Println("IdOK ", productOK.IdOk, " IdBM ", idBM)
				fmt.Println("Name")

				fmt.Println(nameUk == productOK.NameUk)

				fmt.Println(nameRu == productOK.NameRu)
				fmt.Println("Model")

				fmt.Println(modelBM == productOK.Model)
				fmt.Println("Price")
				fmt.Println("BM ", price)
				fmt.Println("OK ", productOK.Price)
				fmt.Println(math.Abs(price-productOK.Price) < 0.01)
				fmt.Println("Brand")
				fmt.Println("BM ", brandId)
				fmt.Println("OK ", productOK.Manufacturer)
				fmt.Println(brandId == productOK.Manufacturer)
				fmt.Println("Category")
				fmt.Println("BM ", catId)
				fmt.Println("OK ", productOK.Category)
				fmt.Println(catId == productOK.Category)
				fmt.Println("Image")
				fmt.Println("BM ", image)
				fmt.Println("OK ", imageOK)
				fmt.Println(image == imageOK)
				fmt.Println("OK Stock ", productOK.Stock)
				fmt.Println(productOK.Stock == 7)
				fmt.Println("OK Quantity ", productOK.Quantity)
				fmt.Println(productOK.Quantity > 0)*/
			}
		}
	}

	fmt.Println("Not checked ", countNotChecked)
	fmt.Println("Not checked Images ", countNotCheckedImg)
	fmt.Println("Not checked Name Uk ", countNotCheckedNameUk)
	fmt.Println("Not checked Name Ru ", countNotCheckedNameRu)
	fmt.Println("Not checked Model ", countNotCheckedModel)
	fmt.Println("Not checked Price ", countNotCheckedPrice)
	fmt.Println("Not checked Brand ", countNotCheckedBrand)
	fmt.Println("Not checked Category ", countNotCheckedCategory)
	fmt.Println("Not checked Stock ", countNotCheckedStock)
	fmt.Println("Not checked Quantity ", countNotCheckedQuantity)
}

func UpdateProdCat(product db.Full) int8 {
	dataUK := product.UK
	pathUK := dataUK["nodes"].(string)
	dataRU := product.RU
	pathRU := dataRU["nodes"].(string)

	client, _ := dbok.NewClient()
	defer client.Close()

	categoryId := GetCategory(client, pathUK, pathRU)
	fmt.Println(categoryId)

	brandUK := dataUK["brand"].(string)
	brandId := GetManufacturer(client, brandUK)
	fmt.Println(brandId)

	filter := bson.D{{"id_bm", product.IdBM}}
	update := bson.D{
		{"$set", bson.D{
			{"category_id", categoryId},
			{"manufacturer_id", brandId},
		},
		},
	}
	_, err := db.ProductFull.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println(err)
	}

	return 0
}

func UpdateProdCatOK(product db.Full, productId int) int8 {
	client, _ := dbok.NewClient()
	defer client.Close()

	categoryId := product.CategoryId
	fmt.Println(categoryId)

	brandId := product.ManufacturerId
	fmt.Println(brandId)

	sql := "INSERT INTO `oc_product_to_category` SET `product_id` = ?, `category_id` = ?"
	_, err := client.Exec(sql, productId, categoryId)
	if err != nil {
		fmt.Println(err)
	}

	sql = "UPDATE `oc_product` SET `manufacturer_id` = ? WHERE `product_id` = ?"
	_, err = client.Exec(sql, brandId, productId)
	if err != nil {
		fmt.Println(err)
	}

	return 0
}

func UpdateProductCat(idBM string) int8 {
	filter := bson.D{{"id_bm", idBM}}
	var product db.Full
	err := db.ProductFull.FindOne(context.TODO(), filter).Decode(&product)

	if err != nil {
		fmt.Println(err)
		return 1
	}

	dataUK := product.UK
	pathUK := dataUK["nodes"].(string)
	dataRU := product.RU
	pathRU := dataRU["nodes"].(string)

	client, _ := dbok.NewClient()
	defer client.Close()

	categoryId := GetCategory(client, pathUK, pathRU)
	fmt.Println(categoryId)

	brandUK := dataUK["brand"].(string)
	brandId := GetManufacturer(client, brandUK)
	fmt.Println(brandId)

	update := bson.D{
		{"$set", bson.D{
			{"category_id", categoryId},
			{"manufacturer_id", brandId},
		},
		},
	}
	_, err = db.ProductFull.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println(err)
	}

	return 0
}

func makeProductMapOK() {
	MapOK = make(map[string]int)

	client, _ := dbok.NewClient()
	defer client.Close()

	var productId int
	var idBM string

	sql := "SELECT `product_id`, `sku` FROM `oc_product` ORDER BY `sku`"
	rows, err := client.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		fmt.Println("From OK.CategoryPath; Count ", count)

		err := rows.Scan(&productId, &idBM)
		if err != nil {
			fmt.Println(err)
			continue
		}

		MapOK[idBM] = productId
	}
}

func makeProductsOK() {
	ProductsOK = make(map[string]OKProduct)

	client, _ := dbok.NewClient()
	defer client.Close()

	var product OKProduct

	sql := "SELECT `p`.`product_id`, `p`.`sku`, `p`.`model`, `pdu`.`name` AS `nameuk`, `pdr`.`name` AS `nameru`, `p`.`price`, `p`.`stock_status_id`, `p`.`manufacturer_id`, `pc`.`category_id`, `p`.`quantity`, `p`.`image` FROM `oc_product` AS `p` LEFT JOIN `oc_product_description` AS `pdu` ON (`p`.`product_id` = `pdu`.`product_id` AND `pdu`.`language_id` = 1) LEFT JOIN `oc_product_description` AS `pdr` ON (`p`.`product_id` = `pdr`.`product_id` AND `pdr`.`language_id` = 3) LEFT JOIN `oc_product_to_category` AS `pc` ON (`p`.`product_id` = `pc`.`product_id`) WHERE `pc`.`category_id` IS NOT NULL ORDER BY `sku`"
	rows, err := client.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		fmt.Println("From OK.Products; Count ", count)

		err := rows.Scan(&product.IdOk, &product.IdBM, &product.Model, &product.NameUk, &product.NameRu, &product.Price, &product.Stock, &product.Manufacturer, &product.Category, &product.Quantity, &product.Image)
		if err != nil {
			fmt.Println(err)
			continue
		}

		ProductsOK[product.IdBM] = product
	}
}

func CheckAbsentCats() {
	makeCategoryMap()
	makeManufacturerMap()

	client, _ := dbok.NewClient()
	defer client.Close()

	var total int
	var idOK int
	var idBM string

	sql := "SELECT COUNT(`p`.`product_id`) FROM `oc_product` AS `p` LEFT JOIN `oc_product_to_category` AS `pc` ON (`p`.`product_id` = `pc`.`product_id`) WHERE `pc`.`category_id` IS NULL ORDER BY `sku`"
	row := client.QueryRow(sql)
	_ = row.Scan(&total)

	sql = "SELECT `p`.`product_id`, `p`.`sku` FROM `oc_product` AS `p` LEFT JOIN `oc_product_to_category` AS `pc` ON (`p`.`product_id` = `pc`.`product_id`) WHERE `pc`.`category_id` IS NULL ORDER BY `sku`"
	rows, err := client.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var productMap map[int]string
	productMap = make(map[int]string)
	count := 0
	for rows.Next() {
		count++
		fmt.Println("From OK.Products without category; Count ", count, " From ", total)

		err := rows.Scan(&idOK, &idBM)
		if err != nil {
			fmt.Println(err)
			continue
		}
		productMap[idOK] = idBM
	}

	var product db.Full
	count = 0
	for idOK, idBM = range productMap {
		count++
		fmt.Println("From products map without category; Count ", count, " From ", total)

		filter := bson.D{{"id_bm", idBM}}
		err = db.ProductFull.FindOne(context.TODO(), filter).Decode(&product)

		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err)
			if AddProductToMongo(idBM) {
				filter := bson.D{{"id_bm", idBM}}
				err = db.ProductFull.FindOne(context.TODO(), filter).Decode(&product)
			}
		}
		if !errors.Is(err, mongo.ErrNoDocuments) {
			dataUK := product.UK
			pathUK := dataUK["nodes"].(string)
			dataRU := product.RU
			pathRU := dataRU["nodes"].(string)

			fmt.Printf("idBM: %#v\n", idBM)
			fmt.Printf("Category: %#v\n", product.CategoryId)
			fmt.Printf("pathUK: %#v\n\n", pathUK)
			fmt.Printf("pathRU: %#v\n\n", pathRU)
			//fmt.Printf("dataUK: %#v\n\n", dataUK)
			//fmt.Printf("dataRU: %#v\n\n", dataRU)
		}
	}

}

func CheckAbsentCats2() {
	client, _ := dbok.NewClient()
	defer client.Close()

	var total int
	var idOK int
	var idBM string

	sql := "SELECT COUNT(`p`.`product_id`) FROM `oc_product` AS `p` LEFT JOIN `oc_product_to_category` AS `pc` ON (`p`.`product_id` = `pc`.`product_id`) WHERE `pc`.`category_id` IS NULL ORDER BY `sku`"
	row := client.QueryRow(sql)
	_ = row.Scan(&total)

	sql = "SELECT `p`.`product_id`, `p`.`sku` FROM `oc_product` AS `p` LEFT JOIN `oc_product_to_category` AS `pc` ON (`p`.`product_id` = `pc`.`product_id`) WHERE `pc`.`category_id` IS NULL ORDER BY `sku`"
	rows, err := client.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var productMap map[int]string
	productMap = make(map[int]string)

	count := 0
	for rows.Next() {
		count++
		fmt.Println("From OK.Products without category; Count ", count, " From ", total)

		err := rows.Scan(&idOK, &idBM)
		if err != nil {
			fmt.Println(err)
			continue
		}

		productMap[idOK] = idBM
	}

	count = 0
	var product db.Full
	for idOK, idBM = range productMap {
		count++
		fmt.Println("From productsMap without category; Count ", count, " From ", total)

		fmt.Println("IdOk: ", idOK, " IdBM: ", idBM)

		filter := bson.D{{"id_bm", idBM}}
		err = db.ProductFull.FindOne(context.TODO(), filter).Decode(&product)

		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println(err)
			continue
		}

		if product.CategoryId > 0 && product.ManufacturerId > 0 {
			fmt.Println(UpdateProdCatOK(product, idOK))
		}
	}
}
