package main

import (
	"bm.parts.server/apiBM"
	"bm.parts.server/config"
	"bm.parts.server/db"
	"bm.parts.server/log"
	"encoding/json"
	"fmt"
	"io"
)

func main() {
	config.Start()
	log.Start()
	defer log.Stop()

	db.Start()
	//defer db.Stop()

	var ret map[string]interface{}
	ret = make(map[string]interface{})

	IdBM := "BFFFC43B104D2532403DE19D114ABD49"

	/*
		GET /product/(string: product_uuid)

		product_uuid – Обов’язковий ID товару
		warehouses – ID складів по яким вертати залишки: &warehouses=816D000C2999A7E611E6EC6B4A1915AF&warehouses=ACF9000C2947F7AE11E28A2B02C4AD32. Значення all верне залишки по всім складам. За замовчуванням, вертає залишки для основного складу.
		currency – ID Валюти для відображення ціни. За замовчуванням, валюта основного договору.
		id_type – Вказує виконувати пошук товару по ID або по коду
		products_as – Формат товарів які повертаються. Можливі варіанти: obj, arr. За замовчуванням, obj.
		q – Пошукова фраза для збереження в історію пошуку.
		save – Прапорець збереження запиту в історію пошуку. За замовчуванням, True. Збереження вимагає заповненого значення параметра q.
		promos – Якщо передати значення full, додасть в відповідь ключ promos_full з інформацією про акції, в яких товар бере участь.
		oe – Формат поля що повертається oe. Можливі варіанти: short - вертає тільки номери, full - повертає додатково конструкційний брендю. За замовчуванням, short.
		output_field – Вказує параметр, який потрібно повернути. Можливі варіанти: all, analogs, cars, oe, buy_with, paired_products, components_of_kit. За замовчуванням, all.
		analogs_available – Не повертати товари-аналоги, яких немає в наявності. Можливі варіанти: 1, 0 (1 - true, 0 - false). За замовчуванням, 0.
	*/
	params := map[string]string{
		"currency":          "A358000C2947F7AE11E23F5617780B16",
		"warehouses":        "all",
		"oe":                "full",
		"output_field":      "all",
		"analogs_available": "0",
	}

	path := "/product/" + IdBM
	retUK, oku := apiBM.Get(path, params, "uk")
	if oku {
		byteValue, _ := io.ReadAll(retUK.Body)
		var result map[string]interface{}
		err := json.Unmarshal(byteValue, &result)
		if err != nil {
			fmt.Println(err)
		}

		ret["UK"] = result["product"]
	}

	fmt.Println(ret)
}
