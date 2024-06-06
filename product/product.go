package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"bm.parts.server/apiBM"
	"bm.parts.server/apiOK"
	"bm.parts.server/db"
	"bm.parts.server/dbok"
	"bm.parts.server/log"
	"bm.parts.server/xtime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var CategoryMutex sync.Mutex
var BrandMutex sync.Mutex

func AddAll() {
	threadsNumber := 12

	filter := bson.D{{"new", 1}}
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
		go addPart(&wg, res)
	} else {
		limit := int(math.Ceil(float64(total / threadsNumber)))
		wg.Add(threadsNumber)
		fmt.Println("wg.Add ", threadsNumber)
		for i := 0; i < threadsNumber; i++ {
			min := i * limit
			max := (i + 1) * limit
			go addPart(&wg, res[min:max])
		}
	}
	wg.Wait()
}

func AddAllToMongo() {
	threadsNumber := 4

	filter := bson.D{{"new", 1}}
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
		go addPartToMongo(&wg, res)
	} else {
		limit := int(math.Ceil(float64(total / threadsNumber)))
		wg.Add(threadsNumber)
		fmt.Println("wg.Add ", threadsNumber)
		for i := 0; i < threadsNumber; i++ {
			min := i * limit
			max := (i + 1) * limit
			go addPartToMongo(&wg, res[min:max])
		}
	}
	wg.Wait()
}

func AddAll3() {
	threadsNumber := 24

	db.ProductEmpty.Drop(context.TODO())

	filter := bson.D{{"in_stock", 1}, {"id_oc", 0}}
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
		go addPart3(&wg, res)
	} else {
		limit := int(math.Ceil(float64(total / threadsNumber)))
		wg.Add(threadsNumber)
		fmt.Println("wg.Add ", threadsNumber)
		for i := 0; i < threadsNumber; i++ {
			min := i * limit
			max := (i + 1) * limit
			go addPart3(&wg, res[min:max])
		}
	}
	wg.Wait()

	emptyCount, errc := db.ProductEmpty.CountDocuments(context.TODO(), bson.D{})
	if errc != nil {
		fmt.Println(errc)
	}
	fmt.Printf("Empty count: %v\n", emptyCount)
}

func AddImg() {
	threadsNumber := 24

	filter := bson.D{{"in_stock", 1}}
	//opts := options.Find().SetLimit(50000).SetSkip(105000)
	//cursor, err := db.ProductShort.Find(context.TODO(), filter, opts)
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
		go addImgs(&wg, res)
	} else {
		limit := int(math.Ceil(float64(total / threadsNumber)))
		wg.Add(threadsNumber)
		fmt.Println("wg.Add ", threadsNumber)
		for i := 0; i < threadsNumber; i++ {
			min := i * limit
			max := (i + 1) * limit
			go addImgs(&wg, res[min:max])
		}
	}
	wg.Wait()
}

func addPart(wg *sync.WaitGroup, rows []db.Short) {
	defer wg.Done()
	defer fmt.Println("wg.Done")

	for i, row := range rows {
		idBM := row.IdBM

		_ = AddProduct(idBM)

		fmt.Println("Add OK; Count ", i, " Row ", row)
	}
}

func addPartToMongo(wg *sync.WaitGroup, rows []db.Short) {
	defer wg.Done()
	defer fmt.Println("wg.Done")

	for i, row := range rows {
		idBM := row.IdBM

		_ = AddProductToMongo(idBM)

		fmt.Println("Add to Mongo; Count ", i, " Row ", row)
	}
}

func addPart3(wg *sync.WaitGroup, rows []db.Short) {
	defer wg.Done()
	defer fmt.Println("wg.Done")

	for i, row := range rows {
		idBM := row.IdBM

		_ = AddProduct3(idBM)

		fmt.Println("Add to Ok; Count ", i, " Row ", row)
	}
}

func addImgs(wg *sync.WaitGroup, rows []db.Short) {
	defer wg.Done()
	defer fmt.Println("wg.Done")

	total := len(rows)
	for i, row := range rows {
		idBM := row.IdBM

		_ = AddImage(idBM)

		fmt.Println("Add Image Site; Count ", i, " From ", total, " Row ", row)
	}
}

func AddProduct(idBM string) bool {
	product, ok := apiBM.GetProduct(idBM)
	if ok {
		dataUK := product["UK"].(map[string]interface{})
		pathUK := dataUK["nodes"].(string)
		dataRU := product["RU"].(map[string]interface{})
		pathRU := dataRU["nodes"].(string)

		brandUK := dataUK["brand"].(string)
		if brandUK == "MOTUL" {
			return false
		}

		client, _ := dbok.NewClient()
		defer client.Close()

		categoryId := GetCategory(client, pathUK, pathRU)

		product["CategoryId"] = categoryId

		brandId := GetManufacturer(client, brandUK)
		product["ManufacturerId"] = brandId

		newProduct := db.Full{
			ID:             primitive.NewObjectID(),
			IdBM:           idBM,
			IdOC:           0,
			IdProm:         0,
			UK:             dataUK,
			RU:             dataRU,
			CategoryId:     categoryId,
			ManufacturerId: brandId,
		}

		result, err := db.ProductFull.InsertOne(context.TODO(), newProduct)
		if err != nil {
			log.MyLog.Println(err)
		}
		fmt.Println(result)

		resp := apiOK.Post("https://dev.dmdshop.com.ua/index.php?route=tool/add", product)
		fmt.Println(resp.Status)

		return true
	}

	return false
}

func AddProductToMongo(idBM string) bool {
	product, ok := apiBM.GetProduct(idBM)
	fmt.Println(ok)
	if ok {
		dataUK := product["UK"].(map[string]interface{})
		pathUK := dataUK["nodes"].(string)
		dataRU := product["RU"].(map[string]interface{})
		pathRU := dataRU["nodes"].(string)

		client, _ := dbok.NewClient()
		defer client.Close()

		categoryId := GetCategory(client, pathUK, pathRU)

		product["CategoryId"] = categoryId

		brandUK := dataUK["brand"].(string)
		brandId := GetManufacturer(client, brandUK)
		product["ManufacturerId"] = brandId

		newProduct := db.Full{
			ID:             primitive.NewObjectID(),
			IdBM:           idBM,
			IdOC:           0,
			IdProm:         0,
			UK:             dataUK,
			RU:             dataRU,
			CategoryId:     categoryId,
			ManufacturerId: brandId,
		}

		result, err := db.ProductFull.InsertOne(context.TODO(), newProduct)
		if err != nil {
			log.MyLog.Println(err)
		}
		fmt.Println(result)

		filter := bson.D{{"id_bm", idBM}}
		update := bson.D{
			{"$set", bson.D{
				{"new", 0},
			},
			},
		}
		_, err = db.ProductShort.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			fmt.Println(err)
		}

		return true
	}

	return false
}

func AddProduct3(idBM string) bool {
	product, ok := GetProductFull(idBM)
	if ok {
		resp := apiOK.PostFull("https://dev.dmdshop.com.ua/index.php?route=tool/add2", product)
		fmt.Println(resp.Status)

		return true
	}

	emptyProd := db.Empty{
		ID:      primitive.NewObjectID(),
		IdBM:    idBM,
		Updated: time.Now().Unix(),
	}
	_, err := db.ProductEmpty.InsertOne(context.TODO(), emptyProd)
	if err != nil {
		log.MyLog.Println(err)
	}

	return false
}

func AddImage(idBM string) bool {
	product, ok := GetProductFull(idBM)
	if ok {
		resp := apiOK.PostFull("https://dev.dmdshop.com.ua/index.php?route=tool/img", product)
		//resp := apiOK.PostFull("http://37.60.251.210/", product) //NKN
		//resp := apiOK.PostFull("https://clickbot.shop/", product)
		//resp := apiOK.PostFull("http://149.102.133.202/", product) // Binance
		fmt.Println(resp.Status)

		return true
	}

	return false
}

func GetProductFull(idBM string) (db.Full, bool) {
	filter := bson.D{{"id_bm", idBM}}
	var product db.Full
	err := db.ProductFull.FindOne(context.TODO(), filter).Decode(&product)

	if errors.Is(err, mongo.ErrNoDocuments) {
		fmt.Println(err)
		if AddProductToMongo(idBM) {
			err = db.ProductFull.FindOne(context.TODO(), filter).Decode(&product)
			if err != nil {
				fmt.Println(err)
				return product, false
			}
		}

		return product, false
	}

	return product, true
}

func GetCategory(client *sql.DB, pathUK string, pathRU string) int {
	CategoryMutex.Lock()
	defer CategoryMutex.Unlock()

	return GetCategoryMute(client, pathUK, pathRU)
}

func GetCategoryMute(client *sql.DB, pathUK string, pathRU string) int {
	nameUK := GetLastName(pathUK)
	nameRU := GetLastName(pathRU)
	pathUK = SanitaizePath(pathUK)
	pathRU = SanitaizePath(pathRU)

	fmt.Println("Path ", pathUK)
	var parentId, top int

	catId, ok := CategoryMap[pathUK]
	if ok {
		return catId
	}

	if strings.Contains(pathUK, "/") {
		top = 0
		parentId = GetCategoryMute(client, CutePath(pathUK), CutePath(pathRU))
	} else {
		top = 1
		parentId = 0
	}

	date := xtime.ToStr(time.Now())
	sql := "INSERT INTO `oc_category` SET `parent_id` = ?, `top` = ?, `column` = 1, sort_order = 0, status = 1, date_added = ?, date_modified = ?"
	res, err := client.Exec(sql, parentId, top, date, date)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	var lastId int64
	lastId, err = res.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}

	catId = int(lastId)

	sql = "INSERT INTO `oc_category_description` SET `category_id` = ?, `language_id` = ?, `name` = ?"
	_, err = client.Exec(sql, catId, 1, nameUK)
	if err != nil {
		fmt.Println(err)
	}
	_, err = client.Exec(sql, catId, 3, nameRU)
	if err != nil {
		fmt.Println(err)
	}

	sql = "INSERT INTO `oc_category_to_store` SET `category_id` = ?, `store_id` = 0"
	_, err = client.Exec(sql, catId)
	if err != nil {
		fmt.Println(err)
	}

	lastLevel := 0
	if parentId > 0 {
		sql = "SELECT `path_id` FROM `oc_category_path` WHERE `category_id` = ? ORDER BY `level`"
		res1, err1 := client.Query(sql, parentId)
		if err1 != nil {
			fmt.Println(err1)
		}
		defer res1.Close()

		for res1.Next() {
			var pathId int

			err = res1.Scan(&pathId)
			if err != nil {
				fmt.Println(err)
				break
			}
			sql = "INSERT INTO `oc_category_path` SET `category_id` = ?, `path_id` = ?, `level` = ?"
			_, err = client.Exec(sql, catId, pathId, lastLevel)
			if err != nil {
				fmt.Println(err)
			}

			lastLevel++
		}
	}
	sql = "INSERT INTO `oc_category_path` SET `category_id` = ?, `path_id` = ?, `level` = ?"
	_, err = client.Exec(sql, catId, catId, lastLevel)
	if err != nil {
		fmt.Println(err)
	}

	sql = "INSERT INTO `oc_category_path_ua` SET `category_id` = ?, `path` = ?"
	_, err = client.Exec(sql, catId, pathUK)
	if err != nil {
		fmt.Println(err)
	}

	CategoryMap[pathUK] = catId

	return catId
}

func SanitaizePath(path string) string {
	path = strings.ToLower(path)
	path = strings.Trim(path, " /")
	path = strings.Replace(path, " /", "/", -1)
	path = strings.Replace(path, "/ ", "/", -1)

	return path
}

func CutePath(path string) string {
	pathArr := strings.Split(path, "/")
	return strings.Join(pathArr[:len(pathArr)-1], "/")
}

func GetLastName(path string) string {
	pathArr := strings.Split(path, "/")

	return pathArr[len(pathArr)-1]
}

func GetManufacturer(client *sql.DB, brand string) int {
	BrandMutex.Lock()
	defer BrandMutex.Unlock()

	brandId, ok := ManufacturerMap[brand]
	if ok {
		return brandId
	}

	sql := "INSERT INTO `oc_manufacturer` SET `name` = ?, `sort_order` = 0, `noindex` = 1"
	res, err := client.Exec(sql, brand)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	var lastId int64
	lastId, err = res.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}

	brandId = int(lastId)

	sql = "INSERT INTO `oc_manufacturer_description` SET `manufacturer_id` = ?, `language_id` = ?"
	_, err = client.Exec(sql, brandId, 1)
	if err != nil {
		fmt.Println(err)
	}
	_, err = client.Exec(sql, brandId, 3)
	if err != nil {
		fmt.Println(err)
	}

	sql = "INSERT INTO `oc_manufacturer_to_store` SET `manufacturer_id` = ?, `store_id` = 0"
	_, err = client.Exec(sql, brandId)
	if err != nil {
		fmt.Println(err)
	}

	ManufacturerMap[brand] = brandId

	return brandId
}
