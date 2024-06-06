package cross

import (
	"archive/zip"
	"bm.parts.server/config"
	"bm.parts.server/db"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"log"
	"net/http"
	"os"
)

func GetCross() {
	res, err := http.Get("https://cdn.bm.parts/promo/b2b_lists/cross_list.zip")
	if err != nil {
		fmt.Println(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		log.Fatal(err)
	}

	// Read all the files from zip archive
	for _, zipFile := range zipReader.File {
		fmt.Println("Reading file:", zipFile.Name)
		dstFile, err := os.OpenFile(config.DirRoot+"/data/"+zipFile.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := zipFile.Open()
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
}

func MakeCrossBase() {
	db.RebuildCross()

	file, err := os.Open(config.DirRoot + "/data/cross_list.csv")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	count := 0
	for i, line := range data {
		count++
		fmt.Println("From Cross; Count ", count)

		if i > 0 {
			newRow := db.Cross{
				ID:          primitive.NewObjectID(),
				Article:     line[0],
				Brand:       line[1],
				Name:        line[2],
				CrossNumber: line[3],
				CrossBrand:  line[4],
			}

			result, err := db.CrossList.InsertOne(context.TODO(), newRow)
			if err != nil {
				log.Println(err)
			}

			log.Println(result)
		}
	}
}
