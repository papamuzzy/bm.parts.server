package autoFilter

import (
	"bm.parts.server/db"
	"bm.parts.server/dbok"
	"context"
	"database/sql"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"strings"
)

type PartsMap map[string]int64
type BrandsMap map[string]int64
type ModelsMap map[int64]map[string]int64
type ModifsMap map[int64]map[string]int64
type PartsToModifMap map[int64]map[int64]bool

type Engine struct {
	EngineType     string `bson:"engine_type"`
	HorsePower     int    `bson:"horse_power"`
	EngineVolume   int    `bson:"engine_volume"`
	Engine         string `bson:"engine"`
	EndRelease     string `bson:"end_release"`
	BodyType       string `bson:"body_type"`
	BeginRelease   string `bson:"begin_release"`
	AmountCylinder int    `bson:"amount_cylinder"`
}

type Model struct {
	Model   string   `bson:"model"`
	Engines []Engine `bson:"engines"`
}

type Car struct {
	Brand  string  `bson:"brand"`
	Models []Model `bson:"models"`
}

var MapParts PartsMap
var MapBrands BrandsMap
var MapModels ModelsMap
var MapModifs ModifsMap
var MapPartToModif PartsToModifMap

var Client *sql.DB

func MakePartsMap() {
	MapParts = make(PartsMap)

	var productId int64
	var sku string

	sql := "SELECT `p`.`product_id`, `p`.`sku` FROM `oc_product` AS `p` LEFT JOIN `oc_auto_link_part` AS `lp` ON (`p`.`product_id` = `lp`.`part_id`) WHERE `lp`.`part_id` IS NULL ORDER BY `product_id`"
	rows, err := Client.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		fmt.Println("Parts Map; Count ", count)

		err := rows.Scan(&productId, &sku)
		if err != nil {
			fmt.Println(err)
		}

		MapParts[sku] = productId
	}
}

func MakeBrandsMap() {
	MapBrands = make(BrandsMap)

	var brandId int64
	var brandName string

	sql := "SELECT `brand_id`, `brand_name` FROM `oc_auto_brands`"
	rows, err := Client.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		fmt.Println("Brands Map; Count ", count)

		err := rows.Scan(&brandId, &brandName)
		if err != nil {
			fmt.Println(err)
		}

		MapBrands[brandName] = brandId
	}
}

func MakeModelsMap() {
	MapModels = make(ModelsMap)

	var brandId int64
	var modelId int64
	var modelName string

	sql := "SELECT `brand_id`, `model_name`, `model_id` FROM `oc_auto_models`"
	rows, err := Client.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		fmt.Println("Models Map; Count ", count)

		err := rows.Scan(&brandId, &modelName, &modelId)
		if err != nil {
			fmt.Println(err)
		}

		_, ok := MapModels[brandId]
		if !ok {
			MapModels[brandId] = make(map[string]int64)
		}

		MapModels[brandId][modelName] = modelId
	}
}

func MakeModifsMap() {
	MapModifs = make(ModifsMap)

	var modelId int64
	var modifId int64
	var modifName string

	sql := "SELECT `model_id`, `modif_name`, `modif_id` FROM `oc_auto_modifications`"
	rows, err := Client.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		fmt.Println("Modifs Map; Count ", count)

		err := rows.Scan(&modelId, &modifName, &modifId)
		if err != nil {
			fmt.Println(err)
		}

		_, ok := MapModifs[modelId]
		if !ok {
			MapModifs[modelId] = make(map[string]int64)
		}

		MapModifs[modelId][modifName] = modifId
	}
}

func MakePartToModifMap() {
	MapPartToModif = make(PartsToModifMap)

	var modifId int64
	var partId int64

	sql := "SELECT `modif_id`, `part_id` FROM `oc_auto_link_part`"
	rows, err := Client.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		fmt.Println("Parts To Modif Map; Count ", count)

		err := rows.Scan(&modifId, &partId)
		if err != nil {
			fmt.Println(err)
		}

		_, ok := MapPartToModif[partId]
		if !ok {
			MapPartToModif[partId] = make(map[int64]bool)
		}

		MapPartToModif[partId][modifId] = true
	}
}

func TrimModifs() {
	MapModifs = make(ModifsMap)

	var modelId int64
	var modifId int64
	var modifName string

	sql := "SELECT `model_id`, `modif_name`, `modif_id` FROM `oc_auto_modifications`"
	rows, err := Client.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		fmt.Println("Trim Modifs Map; Count ", count)

		err := rows.Scan(&modelId, &modifName, &modifId)
		if err != nil {
			fmt.Println(err)
		}

		modifName = strings.Trim(modifName, " ")
		sql = "UPDATE `oc_auto_modifications` SET `modif_name` = ? WHERE `model_id` = ? AND `modif_id` = ?"
		_, err = Client.Exec(sql, modifName, modelId, modifId)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func AddAll() {
	Client, _ = dbok.NewClient()
	defer Client.Close()

	//TrimModifs()

	MakePartsMap()
	MakeBrandsMap()
	MakeModelsMap()
	MakeModifsMap()
	MakePartToModifMap()

	// Get all parts
	limit := int64(100)

	filter := bson.D{{}}
	total, err := db.ProductFull.EstimatedDocumentCount(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
	sliceNumber := int64(math.Ceil(float64(total/limit))) + 1
	fmt.Println("Total ", total, " Slices ", sliceNumber)
	for i := int64(0); i <= sliceNumber; i++ {
		fmt.Println("Slice ", i+1, " From ", sliceNumber)

		skip := i * limit

		opts := options.Find().SetLimit(limit).SetSkip(skip)
		cursor, err := db.ProductFull.Find(context.TODO(), filter, opts)
		if err != nil {
			fmt.Println(err)
			break
		}

		var res []db.Full
		if err = cursor.All(context.TODO(), &res); err != nil {
			fmt.Println(err)
		}

		FilterParts(res)
	}
}

func FilterParts(parts []db.Full) {
	for i, part := range parts {
		fmt.Println("Part ", i, " IdBM ", part.IdBM)

		productId, ok := MapParts[part.IdBM]
		if ok {
			FilterPart(productId, part.UK["cars"].(primitive.A))
		}
	}
}

func FilterPart(productId int64, cars []interface{}) {
	for _, car := range cars {
		brandName := car.(map[string]interface{})["brand"]
		brandId := GetBrandId(brandName.(string))

		fmt.Println("Brand ", brandName, " Id ", brandId)

		for _, model := range car.(map[string]interface{})["models"].(primitive.A) {
			modelName := model.(map[string]interface{})["model"]
			modelId := GetModelId(productId, brandId, modelName.(string))

			fmt.Println("Model ", modelName, " Id ", modelId)

			for _, engine := range model.(map[string]interface{})["engines"].(primitive.A) {
				modifId := GetModifId(modelId, engine.(map[string]interface{}))
				fmt.Println("Engine ", engine, " Id ", modifId)

				if modifId > 0 {
					AddPartToModification(productId, modifId)
				}
			}
		}
	}
}

func GetBrandId(brandName string) int64 {
	brandId, ok := MapBrands[brandName]
	if ok {
		return brandId
	}

	sql := "INSERT INTO `oc_auto_brands` SET `brand_name` = ?, `brand_img` = '', active = 1"
	res, err := Client.Exec(sql, brandName)
	if err != nil {
		fmt.Println(err)
		return int64(0)
	}

	brandId, err = res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return int64(0)
	}

	MapBrands[brandName] = brandId

	fmt.Println("Add Brand ", brandName, " Id ", brandId)

	return brandId
}

func GetModelId(productId int64, brandId int64, modelName string) int64 {
	modelId, ok := MapModels[brandId][modelName]
	if ok {
		return modelId
	}

	sql := "INSERT INTO `oc_auto_models` SET `model_name` = ?, `brand_id` = ?, active = 1, `model_start` = '', `model_end` = '', `model_img` = ''"
	res, err := Client.Exec(sql, modelName, brandId)
	if err != nil {
		fmt.Println(err)
		return int64(0)
	}

	modelId, err = res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return int64(0)
	}

	_, ok = MapModels[brandId]
	if !ok {
		MapModels[brandId] = make(map[string]int64)
	}

	MapModels[brandId][modelName] = modelId

	modifId := GetModifEmptyId(modelId)
	AddPartToModification(productId, modifId)

	fmt.Println("Add Model ", modelName, " Brand ", brandId)

	return modelId
}

func GetModifEmptyId(modelId int64) int64 {
	var modifId int64

	sql := "INSERT INTO `oc_auto_modifications` SET `model_id` = ?, `modif_name` = 'all', `modif_dvs` = '-', `modif_hp` = '-', `modif_start` = '-', `modif_end` = '-', `modif_body` = '-', `active` = 1"
	res, err := Client.Exec(sql, modelId)
	if err != nil {
		fmt.Println(err)
		return int64(0)
	}

	modifId, err = res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return int64(0)
	}

	fmt.Println("Add empty modif to model ", modelId)

	return modifId
}

func GetModifId(modelId int64, engine map[string]interface{}) int64 {
	name := engine["engine_type"].(string) + " "
	name += engine["engine"].(string) + " "
	name += engine["body_type"].(string)
	name = strings.Trim(name, " ")

	modifId, ok := MapModifs[modelId][name]
	if ok {
		return modifId
	}

	start := ""
	end := ""
	body := ""
	if engine["begin_release"] != nil {
		start = engine["begin_release"].(string)[:4]
	}
	if engine["end_release"] != nil {
		end = engine["end_release"].(string)[:4]
	}
	if engine["body_type"] != nil {
		body = engine["body_type"].(string)
	}

	sql := "INSERT INTO `oc_auto_modifications` SET `model_id` = ?, `modif_name` = ?, `modif_dvs` = '', `modif_hp` = '', `modif_start` = ?, `modif_end` = ?, `modif_body` = ?, `active` = 1"
	fmt.Println("sql ", sql, "\nname >", name, "< modelId >", modelId, "< start >", start, "< end >", end, "< body >", body, "<")
	res, err := Client.Exec(sql)
	if err != nil {
		fmt.Println(err)
		return int64(0)
	}

	modifId, err = res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return int64(0)
	}

	_, ok = MapModifs[modelId]
	if !ok {
		MapModifs[modelId] = make(map[string]int64)
	}
	MapModifs[modelId][name] = modifId

	fmt.Println("Add modif ", name, " to Model ", modelId)

	return modifId
}

func AddPartToModification(productId int64, modifId int64) {
	_, ok := MapPartToModif[productId][modifId]
	if !ok {
		sql := "INSERT INTO `oc_auto_link_part` SET `modif_id` = ?, `part_id` = ?"
		_, err := Client.Exec(sql, modifId, productId)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Product NOT added!!!")
			return
		}

		_, ok = MapPartToModif[productId]
		if !ok {
			MapPartToModif[productId] = make(map[int64]bool)
		}
		MapPartToModif[productId][modifId] = true

		fmt.Println("Add product ", productId, "to modification ", modifId)
		return
	}

	fmt.Println("Product added")
}
