package dbok

import (
	"bm.parts.server/xtime"
	"database/sql"
	"fmt"
	"time"
)

func ProductUpdate(client *sql.DB, productId int64, price float64, inStock int8) {
	var quantity int8
	if inStock == 0 {
		quantity = 0
	} else {
		quantity = 100
	}
	var stockStatus int8
	if inStock == 0 {
		stockStatus = 5
	} else if inStock == 1 {
		stockStatus = 7
	} else if inStock == 2 {
		stockStatus = 6
	}
	dateModified := xtime.ToStr(time.Now())

	res, err := client.Exec("UPDATE `oc_product` SET `quantity` = ?, `stock_status_id` = ?, `price` = ?, `date_modified` = ? WHERE `product_id` = ?", quantity, stockStatus, price, dateModified, productId)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("res", res)
}

func DisableAllProducts(client *sql.DB) {
	dateModified := xtime.ToStr(time.Now())
	_, err := client.Exec("UPDATE `oc_product` SET `quantity` = ?, `stock_status_id` = ?, `date_modified` = ? WHERE 1", 0, 5, dateModified)

	if err != nil {
		fmt.Println(err)
	}
}
