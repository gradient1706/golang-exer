package dao

import (
	fmt "fmt"
	"log"
	api_package "models"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"

	// "gopkg.in/mgo.v2/bson"
	bson_api "github.com/globalsign/mgo/bson"
)

type ProductDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "product"
)

// Establish a connection to database
func (m *ProductDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Find list of Products
func (m *ProductDAO) FindAll() ([]api_package.Product, error) {
	var Products []api_package.Product
	err := db.C(COLLECTION).Find(bson_api.M{}).All(&Products)
	return Products, err
}

// PurchaseProductEndPoint
func (m *ProductDAO) PurchaseProductDao(products []api_package.Product, quantities []int) error {
	var ops []txn.Op
	runner := txn.NewRunner(db.C(COLLECTION))

	for index, product := range products {
		fmt.Println(product.ID)
		fmt.Println(product.ProductID)
		fmt.Println(product.Quantity)
		op := txn.Op{
			C:      "product",
			Id:     product.ID,
			Assert: bson_api.M{"quantity": bson_api.M{"$gte": quantities[index]}},
			Update: bson.M{"$inc": bson.M{"quantity": -quantities[index]}},
		}

		ops = append(ops, op)
	}
	id := bson.NewObjectId() // Optional
	err := runner.Run(ops, id, nil)
	if err != nil {
		panic(err)
	}
	return err
}

// Find a product by its productID
func (m *ProductDAO) FindOne(productID int) (api_package.Product, error) {

	var product api_package.Product
	err := db.C(COLLECTION).Find(bson_api.M{"productID": productID}).One(&product)
	return product, err
}

// Find a product by its id
func (m *ProductDAO) FindById(id string) (api_package.Product, error) {
	var movie api_package.Product
	err := db.C(COLLECTION).FindId(bson_api.ObjectIdHex(id)).One(&movie)
	return movie, err
}

// Insert a product into database
func (m *ProductDAO) Insert(product api_package.Product) error {
	err := db.C(COLLECTION).Insert(&product)
	return err
}

// Delete an existing product
func (m *ProductDAO) Delete(product api_package.Product) error {
	err := db.C(COLLECTION).Remove(&product)
	return err
}

// Update an existing product
func (m *ProductDAO) Update(product api_package.Product) error {
	err := db.C(COLLECTION).UpdateId(product.ID, &product)
	return err
}
