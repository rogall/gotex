package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Product struct {
	gorm.Model
	Title  string
	Description string
}

func allProducts(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var products []Product
	db.Find(&products)
	fmt.Println("{}", products)

	json.NewEncoder(w).Encode(products)
}

func newProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New Product Endpoint Hit")

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	title := vars["title"]
	description := vars["description"]

	fmt.Println(title)
	fmt.Println(description)

	db.Create(&Product{Title: title, Description: description})
	fmt.Fprintf(w, "New Product Successfully Created")
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	title := vars["title"]

	var product Product
	db.Where("title = ?", title).Find(&product)
	db.Delete(&product)

	fmt.Fprintf(w, "Successfully Deleted Product")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/products", allProducts).Methods("GET")
	myRouter.HandleFunc("/product/{name}", deleteProduct).Methods("DELETE")	
	myRouter.HandleFunc("/product/{title}/{description}", newProduct).Methods("POST")
	log.Fatal(http.ListenAndServe("localhost:8081", myRouter))
}

func initialMigration() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Product{})
}

func main() {
	fmt.Println("Go ORM Tutorial")

	initialMigration()
	// Handle Subsequent requests
	handleRequests()
}
