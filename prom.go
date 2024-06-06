package main

import (
	"bm.parts.server/apiProm"
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

	path := "/orders/list"
	params := make(map[string]string)
	params["limit"] = "100"

	res, ok := apiProm.Get(path, params)
	if ok {
		fmt.Printf("%#v", res)

		byteValue, _ := io.ReadAll(res.Body)
		var result apiProm.Orders
		err := json.Unmarshal(byteValue, &result)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("\nResult %#v\n", result)
	}

}
