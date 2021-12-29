package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"codetest-docker/app/controllers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	// github actionsではdbがコンテナではないので、"172.18.0.1"を環境変数で指定しなおす
	db, err := sql.Open("mysql", "root@tcp("+os.Getenv("DB_HOST")+")/codetest")
	if err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()
	controller := controllers.Controller{}
	router.HandleFunc("/users", controller.GetUsers(db)).Methods("GET")
	router.HandleFunc("/transactions", controller.GetTransactions(db)).Methods("GET")
	router.HandleFunc("/transactions", controller.AddTransaction(db)).Methods("POST")

	log.Println("Server up on port 8888...")
	log.Fatal(http.ListenAndServe(":8888", router))
}
