package db

import (
	"bm.parts.server/config"
	"bm.parts.server/log"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var DataBase *mongo.Database
var ProductShort *mongo.Collection
var ProductFull *mongo.Collection
var ProductEmpty *mongo.Collection
var CrossList *mongo.Collection

func Start() {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.MongoUri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	Client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.MyLog.Println(err)
	}

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := Client.Database(config.MongoBase).RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		log.MyLog.Println(err)
	}
	log.MyLog.Println("Pinged your deployment. You successfully connected to MongoDB!")

	DataBase = Client.Database(config.MongoBase)
	log.MyLog.Println("DataBase " + config.MongoBase + " successfully connected to MongoDB!")
	ProductShort = DataBase.Collection(config.MongoCollection)
	log.MyLog.Println("Collection " + config.MongoCollection + " successfully connected to MongoDB!")
	ProductFull = DataBase.Collection(config.MongoCollection2)
	log.MyLog.Println("Collection " + config.MongoCollection2 + " successfully connected to MongoDB!")
	ProductEmpty = DataBase.Collection("product_empty")
	log.MyLog.Println("Collection product_empty successfully connected to MongoDB!")
	CrossList = DataBase.Collection("cross")
	log.MyLog.Println("Collection cross successfully connected to MongoDB!")
}

func CopyFull() {
	fullTmp := DataBase.Collection("product_full_tmp")

	cursor, err := ProductFull.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	count := 0
	for cursor.Next(context.TODO()) {
		var result Full
		if err := cursor.Decode(&result); err != nil {
			fmt.Println(err)
		}

		res, err := fullTmp.InsertOne(context.TODO(), result)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)

		count++
		fmt.Println("Count ", count)
	}
	if err := cursor.Err(); err != nil {
		panic(err)
	}
}

func CopyToFull() {
	fullTmp := DataBase.Collection("product_full_tmp")

	cursor, err := fullTmp.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	count := 0
	for cursor.Next(context.TODO()) {
		var result Full
		if err := cursor.Decode(&result); err != nil {
			fmt.Println(err)
		}

		res, err := ProductFull.InsertOne(context.TODO(), result)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)

		count++
		fmt.Println("Count ", count)
	}
	if err := cursor.Err(); err != nil {
		panic(err)
	}
}

func DropTmp() {
	fullTmp := DataBase.Collection("product_full_tmp")
	if err := fullTmp.Drop(context.TODO()); err != nil {
		fmt.Println(err)
	}
}

func RebuildProductShort() {
	if err := ProductShort.Drop(context.TODO()); err != nil {
		fmt.Println(err)
	}
	indexes := []mongo.IndexModel{{Keys: bson.D{{"id_bm", 1}}}, {Keys: bson.D{{"id_oc", 1}}}, {Keys: bson.D{{"id_prom", 1}}}}
	names, err := ProductShort.Indexes().CreateMany(context.TODO(), indexes)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Names of Index Created: %v", names)
}

func RebuildProductFull() {
	if err := ProductFull.Drop(context.TODO()); err != nil {
		fmt.Println(err)
	}
	indexes := []mongo.IndexModel{{Keys: bson.D{{"id_bm", 1}}}}
	names, err := ProductFull.Indexes().CreateMany(context.TODO(), indexes)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Names of Index Created: %v", names)
}

func RebuildCross() {
	if err := CrossList.Drop(context.TODO()); err != nil {
		fmt.Println(err)
	}
	indexes := []mongo.IndexModel{{Keys: bson.D{{"article", 1}}}, {Keys: bson.D{{"brand", 1}}}}
	names, err := CrossList.Indexes().CreateMany(context.TODO(), indexes)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Names of Index Created: %v", names)
}

func Stop() {
	if err := Client.Disconnect(context.TODO()); err != nil {
		log.MyLog.Println(err)
	}
}

func GetProductShort(idBM string) (Short, bool) {
	var product Short
	ok := true

	filter := bson.D{{"IdBM", idBM}}
	err := ProductShort.FindOne(context.TODO(), filter).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ok = false
		}
		log.MyLog.Println(err)
	}

	return product, ok
}
