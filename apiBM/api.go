package apiBM

import (
	"bm.parts.server/config"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// GetProduct
//
// GET /product/(string: product_uuid)
//
// product_uuid – Обов’язковий ID товару
// warehouses – ID складів по яким вертати залишки: &warehouses=816D000C2999A7E611E6EC6B4A1915AF&warehouses=ACF9000C2947F7AE11E28A2B02C4AD32.
//
//	Значення all верне залишки по всім складам. За замовчуванням, вертає залишки для основного складу.
//
// currency – ID Валюти для відображення ціни. За замовчуванням, валюта основного договору.
// id_type – Вказує виконувати пошук товару по ID або по коду
// products_as – Формат товарів які повертаються. Можливі варіанти: obj, arr. За замовчуванням, obj.
// q – Пошукова фраза для збереження в історію пошуку.
// save – Прапорець збереження запиту в історію пошуку. За замовчуванням, True. Збереження вимагає заповненого значення параметра q.
// promos – Якщо передати значення full, додасть в відповідь ключ promos_full з інформацією про акції, в яких товар бере участь.
// oe – Формат поля що повертається oe. Можливі варіанти: short - вертає тільки номери, full - повертає додатково конструкційний бренд. За замовчуванням, short.
// output_field – Вказує параметр, який потрібно повернути. Можливі варіанти: all, analogs, cars, oe, buy_with, paired_products, components_of_kit. За замовчуванням, all.
// analogs_available – Не повертати товари-аналоги, яких немає в наявності. Можливі варіанти: 1, 0 (1 - true, 0 - false). За замовчуванням, 0.
func GetProduct(idBM string) (map[string]interface{}, bool) {
	ret := make(map[string]interface{})

	langs := map[string]string{
		"uk": "UK",
		"ru": "RU",
	}

	params := map[string]string{
		"currency":          "A358000C2947F7AE11E23F5617780B16",
		"warehouses":        "all",
		"oe":                "full",
		"output_field":      "all",
		"analogs_available": "0",
	}

	path := "/product/" + idBM

	for lng, code := range langs {
		res, ok := Get(path, params, lng)
		if ok {
			byteValue, _ := io.ReadAll(res.Body)
			var result map[string]interface{}
			err := json.Unmarshal(byteValue, &result)
			if err != nil {
				return nil, false
			}

			ret[code] = result["product"]
		} else {
			return nil, false
		}
	}

	return ret, true
}

func GetProductShort(idBM string) (map[string]interface{}, bool) {
	ret := make(map[string]interface{})

	params := map[string]string{
		"currency":          "A358000C2947F7AE11E23F5617780B16",
		"warehouses":        "all",
		"oe":                "full",
		"output_field":      "all",
		"analogs_available": "0",
	}

	path := "/product/" + idBM
	res, ok := Get(path, params, "uk")
	if ok {
		byteValue, _ := io.ReadAll(res.Body)
		var result map[string]interface{}
		err := json.Unmarshal(byteValue, &result)
		if err != nil {
			return nil, false
		}

		ret = result["product"].(map[string]interface{})
	} else {
		return nil, false
	}

	return ret, true
}

// GetPrice
// /prices/list
func GetPrice() {
	ret := Post("/prices/list", map[string]interface{}{"currency": "A358000C2947F7AE11E23F5617780B16"})

	file, err := os.Create(config.DirRoot + "/data/BMPrice.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	csv, err := io.ReadAll(ret.Body)
	if err != nil {
		fmt.Println(err)
	}

	_, err = fmt.Fprint(file, string(csv))
}
