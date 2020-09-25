package main

import (
	api_config "github.com/gradient1706/golang-exer/config"
	api_dao "github.com/gradient1706/golang-exer/dao"
	"encoding/json"
	fmt "fmt"
	"log"
	api_package "github.com/gradient1706/golang-exer/models"
	"net/http"
	"github.com/gorilla/mux"
	// "time"
	// "context"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// bson_api "github.com/globalsign/mgo/bson"
    // "go.mongodb.org/mongo-driver/mongo/readpref"
)

var config = api_config.Config{}
var dao = api_dao.ProductDAO{}

// // GET list of Products
func AllProductsEndPoint(w http.ResponseWriter, r *http.Request) {

	Products, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, Products)
}

// GET a product by its ID
func FindProductEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	product, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	respondWithJson(w, http.StatusOK, product)
}

func PurchaseProductEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var reqProducts []api_package.Product

	if err := json.NewDecoder(r.Body).Decode(&reqProducts); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	fmt.Print(reqProducts)
	var products []api_package.Product
	var quantities []int
	for _, reqProduct := range reqProducts {
		var reqQuantity = reqProduct.Quantity
		resproduct, err := dao.FindOne(reqProduct.ProductID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
			return
		}
		var stockQuantity = resproduct.Quantity

		if stockQuantity < reqQuantity {
			respondWithJson(w, http.StatusOK, map[string]string{"successful": "false"})
			return
		}
		//fmt.Print(resproduct.Quantity)
		products = append(products, resproduct)
		quantities = append(quantities, reqProduct.Quantity)
		// resproduct.Quantity = stockQuantity - reqQuantity
		// if err := dao.Update(resproduct); err != nil {
		// 	respondWithError(w, http.StatusInternalServerError, err.Error())
		// 	return
		// }
	}

	err := dao.PurchaseProductDao(products, quantities)
	if err != nil {
		respondWithJson(w, http.StatusOK, map[string]string{"successful": "false"})
		return
	}

	respondWithJson(w, http.StatusOK, map[string]string{"successful": "true"})
}

// POST a new product
func CreateProductEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var product api_package.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	product.ID = primitive.NewObjectID()
	if err := dao.Insert(product); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, product)
}

// // PUT update an existing product
func UpdateProductEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var product api_package.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(product); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// // DELETE an existing product
func DeleteProductEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var product api_package.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Delete(product); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// // Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

func main() {
	//client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://user:6huJmxlyu5fegs3G@cluster0.5sgx9.mongodb.net/Products_db?retryWrites=true&w=majority"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// err = client.Connect(ctx)
	// if err != nil{
	// 	log.Fatal(err)
	// }
	// defer  client.Disconnect(ctx)
	
	// productDatabase := client.Database("Products_db")
	//productCollection := productDatabase.Collection("product")

	// Post one
	// productResult, err := productCollection.InsertOne(ctx, bson.D{
	// 	{Key: "productID", Value: 2},
	// 	{Key: "quantity", Value: 5},
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Inserted %v documents into product collection!\n", productResult.InsertedID)

	

	//Get one
	// var product bson.M
	// if err = productCollection.FindOne(ctx, bson.M{}).Decode(&product); err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(product)

	// Get with a filter
	// filterCursor, err := productCollection.Find(ctx, bson.M{"productID": 1})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// var productsFiltered []bson.M
	// if err = filterCursor.All(ctx, &productsFiltered); err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(productsFiltered)

	

	// Delete
	// objID, err := primitive.ObjectIDFromHex("5f4a9cfa871026f9d1334508")
	// result, err := productCollection.DeleteOne(ctx, bson.M{"_id": objID})
	
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("DeleteOne removed %v document(s)\n", result.DeletedCount)

	// Update but keep id
	// result, err := productCollection.ReplaceOne(
	// 	ctx,
	// 	bson.M{"productID": 1},
	// 	bson.M{
	// 		"productID": 1,
	// 		"quantity": 10,
	// 	},
	// )
	// fmt.Printf("Replaced %v Documents!\n", result.ModifiedCount)

	// Get all
	// cursor, err := productCollection.Find(ctx, bson.M{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// var products []bson.M
	// if err = cursor.All(ctx, &products); err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(products)

	// Map with object
	// var products []api_package.Product
	// cursor, err := productCollection.Find(ctx, bson.M{"productID": 1})
	// if err != nil {
	// 	panic(err)
	// }
	// if err = cursor.All(ctx, &products); err != nil {
	// 	panic(err)
	// }
	// fmt.Println(products)

	r := mux.NewRouter()
	r.HandleFunc("/products", AllProductsEndPoint).Methods("GET")
	r.HandleFunc("/products", CreateProductEndPoint).Methods("POST")
	r.HandleFunc("/purchase", PurchaseProductEndPoint).Methods("POST")
	r.HandleFunc("/products", UpdateProductEndPoint).Methods("PUT")
	r.HandleFunc("/products", DeleteProductEndPoint).Methods("DELETE")
	r.HandleFunc("/products/{id}", FindProductEndpoint).Methods("GET")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
