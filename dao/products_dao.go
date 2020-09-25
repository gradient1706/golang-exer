package dao

import (
	fmt "fmt"
	"log"
	api_package "github.com/gradient1706/golang-exer/models"
	// api_config "github.com/gradient1706/golang-exer/config"
	// api_dao "github.com/gradient1706/golang-exer/dao"
	// "encoding/json"
	// "time"
	"context"
	"go.mongodb.org/mongo-driver/mongo/writeconcern" 
	"go.mongodb.org/mongo-driver/mongo/readconcern" 
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// bson_api "github.com/globalsign/mgo/bson"
    // "go.mongodb.org/mongo-driver/mongo/readpref"
)

type ProductDAO struct {
	Server   string
	Database string
}

var cl *mongo.Client
var db *mongo.Database

const (
	COLLECTION = "product"
)

func (m *ProductDAO) Connect() {
	// session, err := mgo.Dial(m.Server)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// db = session.DB(m.Database)

	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://user:6huJmxlyu5fegs3G@cluster0.5sgx9.mongodb.net/Products_db?retryWrites=true&w=majority"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer  client.Disconnect(context.TODO())
	
	// db = client.Database("Products_db")


	clientOptions := options.Client().ApplyURI("mongodb+srv://user:6huJmxlyu5fegs3G@cluster0.5sgx9.mongodb.net/Products_db?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	cl = client
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	db = client.Database("Products_db")
}


// Find list of Products
func (m *ProductDAO) FindAll() ([]api_package.Product, error) {
	// var Products []api_package.Product
	// err := db.C(COLLECTION).Find(bson_api.M{}).All(&Products)
	// return Products, err

	productCollection := db.Collection("product")
	var products []api_package.Product
	cursor, err := productCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		panic(err)
	}
	if err = cursor.All(context.TODO(), &products); err != nil {
		panic(err)
	}
	fmt.Println(products)
	return products, err
}

// PurchaseProductEndPoint
func (m *ProductDAO) PurchaseProductDao(products []api_package.Product, quantities []int) error {
	productCollection := db.Collection("product")
	wc := writeconcern.New(writeconcern.WMajority())
    rc := readconcern.Snapshot()
    txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

    session, err := cl.StartSession()
    if err != nil {
        panic(err)
    }
    defer session.EndSession(context.Background())

    callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		
		var result mongo.UpdateResult
		for index, element := range products {
			filter := bson.D{{"productID", element.ProductID}}

			update := bson.D{
				{"$inc", bson.D{
					{"quantity", -quantities[index]},
				}},
			}

			_, err := productCollection.UpdateOne(
				sessionContext,
				filter,
				update,
			)
			if err != nil {
				return nil, err
			}
		}

        return result, err
    }

    _, err = session.WithTransaction(context.Background(), callback, txnOpts)
    if err != nil {
        panic(err)
	}
	return err


}

// Find a product by its productID
func (m *ProductDAO) FindOne(productID int) (api_package.Product, error) {

	// var product api_package.Product
	// err := db.C(COLLECTION).Find(bson_api.M{"productID": productID}).One(&product)
	// return product, err

	var product api_package.Product
	productCollection := db.Collection("product")
	filter := bson.D{{"productID", productID}}
	err := productCollection.FindOne(context.TODO(), filter).Decode(&product)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", product)
	return product, err
}

// Find a product by its id
func (m *ProductDAO) FindById(id string) (api_package.Product, error) {
	// var movie api_package.Product
	// err := db.C(COLLECTION).FindId(bson_api.ObjectIdHex(id)).One(&movie)
	// return movie, err

	var product api_package.Product
	productCollection := db.Collection("product")
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objID}}
	err := productCollection.FindOne(context.TODO(), filter).Decode(&product)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", product)
	return product, err
}

// Insert a product into database
func (m *ProductDAO) Insert(product api_package.Product) error {

	productCollection := db.Collection("product")
	insertResult, err := productCollection.InsertOne(context.TODO(), product)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return err
}

// Delete an existing product
func (m *ProductDAO) Delete(product api_package.Product) error {
	
	productCollection := db.Collection("product")
	filter := bson.D{{"_id", product.ID}}
	deleteResult, err := productCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	return err
}

// Update an existing product
func (m *ProductDAO) Update(product api_package.Product) error {

	productCollection := db.Collection("product")
	result, err := productCollection.ReplaceOne(
		context.TODO(),
		bson.M{"_id": product.ID},
		bson.M{
			"productID": product.ProductID,
			"quantity": product.Quantity,
		},
	)
	fmt.Printf("Replaced %v Documents!\n", result.ModifiedCount)
	return err
}
