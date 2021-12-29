package main

import (
	"database/sql"
	"log"
	"net/http"

	"codetest-docker/app/controllers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	db, err := sql.Open("mysql", "root@tcp(db:3306)/codetest")
	if err != nil {
		log.Fatal(err)
	}
	controller := controllers.Controller{}
	router := mux.NewRouter()
	router.HandleFunc("/api/users", controller.GetUsers(db)).Methods("GET")
	log.Println("Server up on port 8888...")
	log.Fatal(http.ListenAndServe(":8888", router))
}
