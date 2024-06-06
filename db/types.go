package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type Short struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	IdBM    string             `bson:"id_bm" validate:"required"`
	IdOC    int64              `bson:"id_oc"`
	IdProm  int64              `bson:"id_prom"`
	Price   float64            `bson:"price"`
	InStock int8               `bson:"in_stock"` // 0 -- not in price; 1 -- In price; 2 -- Expected
	Updated int64              `bson:"updated"`
	New     int8               `bson:"new"`
}

type Full struct {
	ID             primitive.ObjectID     `bson:"_id,omitempty"`
	IdBM           string                 `bson:"id_bm" validate:"required"`
	IdOC           int64                  `bson:"id_oc"`
	IdProm         int64                  `bson:"id_prom"`
	UK             map[string]interface{} `bson:"uk"`
	RU             map[string]interface{} `bson:"ru"`
	CategoryId     int                    `bson:"category_id"`
	ManufacturerId int                    `bson:"manufacturer_id"`
}

type Empty struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	IdBM    string             `bson:"id_bm" validate:"required"`
	Updated int64              `bson:"updated"`
}

type Cross struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Article     string             `bson:"article"`
	Brand       string             `bson:"brand"`
	Name        string             `bson:"name"`
	CrossNumber string             `bson:"cross_number"`
	CrossBrand  string             `bson:"cross_brand"`
}
