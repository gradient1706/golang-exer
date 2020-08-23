package main

import (
	api_config "config"
	api_dao "dao"
	"encoding/json"
	fmt "fmt"
	"log"
	api_package "models"
	"net/http"

	bson_api "github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

var config = api_config.Config{}
var dao = api_dao.ProductDAO{}

// GET list of Products
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
	product.ID = bson_api.NewObjectId()
	if err := dao.Insert(product); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, product)
}

// PUT update an existing product
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

// DELETE an existing product
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

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

// Define HTTP request routes
func main() {
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
