package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var IsDebug bool
var IsTest bool
var DirRoot string
var ServerLog string
var UpdateLog string

// MongoUri Mongo params
var MongoUri string
var MongoBase string
var MongoCollection string
var MongoCollection2 string

// OpenCart params
var OcServer string
var OcPort string
var OcDriver string
var OcUser string
var OcPassw string
var OcDataBase string

// Prom params
var PromUrl string
var PromToken string

// BM Parts params
var BmUrl string
var BmToken string
var BmUserAgent string

func Start() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	IsDebug = os.Getenv("DEBUG") == "1"
	IsTest = os.Getenv("TEST") == "1"

	ServerLog = os.Getenv("SERVER_LOG")
	UpdateLog = os.Getenv("UPDATE_LOG")

	MongoUri = os.Getenv("MONGO_URI")
	MongoBase = os.Getenv("MONGO_BASE")
	MongoCollection = os.Getenv("MONGO_COLLECTION")
	MongoCollection2 = os.Getenv("MONGO_COLLECTION2")

	OcServer = os.Getenv("OC_SERVER")
	OcPort = os.Getenv("OC_PORT")
	OcDriver = os.Getenv("OC_DRIVER")
	OcUser = os.Getenv("OC_USER")
	OcPassw = os.Getenv("OC_PASSW")
	OcDataBase = os.Getenv("OC_DATABASE")

	PromUrl = os.Getenv("PROM_URL")
	PromToken = os.Getenv("PROM_TOKEN")

	BmUrl = os.Getenv("BM_URL")
	BmToken = os.Getenv("BM_TOKEN")
	BmUserAgent = os.Getenv("BM_USER_AGENT")

	/*
		Port, err = strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			log.Fatalf("Port error. Err: %s", err)
		}

		NodeNum, err = strconv.Atoi(os.Getenv("NODE_NUMBER"))
		if err != nil {
			NodeNum = 1
		}
		minProfit, _ = strconv.ParseFloat(os.Getenv("MIN_PROFIT"), 64)
		maxLevel, _ = strconv.Atoi(os.Getenv("MAX_LEVEL"))
		fiat = strings.Split(os.Getenv("FIAT"), ",")
	*/

	getRoot()
}

func getRoot() {
	DirRoot, _ = os.Getwd()
	fmt.Println(DirRoot)
}
