package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
    ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    ProductID int `bson:"productID" json:"productID"`
    Quantity int `bson:"quantity" json:"quantity"`
}
