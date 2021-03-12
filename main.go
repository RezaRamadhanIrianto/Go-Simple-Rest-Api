package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shopspring/decimal"
)

var db *gorm.DB
var err error

type Product struct{
	ID int `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	Price decimal.Decimal `json:"price" sql:"type:decimal(16,2)"`
}

type Result struct{
	Code int `json:"code"`
	Data interface{} `json:"data"`
	Message string `json:"message"`
}

func main() {
	db, err = gorm.Open("mysql", "root:root@(127.0.0.1:3306)/rest_go")

	if err != nil {
		log.Println("Connection Failed", err)
	} else {
		log.Println("Connection Establised")
	}

	db.AutoMigrate(&Product{})

	handleRequests()
}


func handleRequests(){
	log.Println("Start the development at https://127.0.0.1:9999")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/api/products", createProduct).Methods("POST")
	myRouter.HandleFunc("/api/products", getProducts).Methods("GET")
	myRouter.HandleFunc("/api/products/{id}", getProduct).Methods("GET")
	myRouter.HandleFunc("/api/products/{id}", updateProduct).Methods("PUT")
	myRouter.HandleFunc("/api/products/{id}", deleteProduct).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9999", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Welcome ")
}

func createProduct(w http.ResponseWriter, r *http.Request){
	payloads, _ := ioutil.ReadAll(r.Body)
	var product Product
	json.Unmarshal(payloads, &product)
	db.Create(&product)


	res := Result{Code: 200, Data: product, Message: "Success create product"}
		
	result, err := json.Marshal(res)

	if(err != nil){
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getProducts(w http.ResponseWriter, r *http.Request){
	products := []Product{}
	db.Find(&products)

	res := Result{Code: 200, Data: products, Message: "Success get products"}
	result, err := json.Marshal(res)

	if(err != nil){
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getProduct(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	productId := vars["id"]


	var product Product
	db.First(&product, productId)

	res := Result{Code: 200, Data: product, Message: "Success get products"}
	result, err := json.Marshal(res)

	if(err != nil){
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func updateProduct(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	productId := vars["id"]

	payloads, _ := ioutil.ReadAll(r.Body)
	var productUpdate Product

	json.Unmarshal(payloads, &productUpdate)
	
	var product Product
	db.First(&product, productId)
	db.Model(&product).Updates(productUpdate)


	res := Result{Code: 200, Data: product, Message: "Success Update product"}
		
	result, err := json.Marshal(res)

	if(err != nil){
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func deleteProduct(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	productId := vars["id"]

	var product Product
	db.First(&product, productId)
	db.Delete(&product)
	res := Result{Code: 200, Message: "Success Delete product"}
		
	result, err := json.Marshal(res)

	if(err != nil){
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}