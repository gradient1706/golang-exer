package models

import (
	// "gopkg.in/mgo.v2/bson"
	"github.com/globalsign/mgo/bson" 
)

type Product struct {
    ID bson.ObjectId `bson:"_id" json:"id"`
    ProductID int `bson:"productID" json:"productID"`
    Quantity int `bson:"quantity" json:"quantity"`
}
