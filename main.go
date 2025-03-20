package main

import (
	"log"
	"net/http"

	"todo-app/database"
	"todo-app/handlers"

	"github.com/gorilla/mux"
)

func main() {

	database.InitDB()
	router := mux.NewRouter()

	router.HandleFunc("/tasks", handlers.PostTaskHandler).Methods("POST")
	router.HandleFunc("/tasks/{id}", handlers.PutTaskHandler).Methods("PUT")
	router.HandleFunc("/tasks", handlers.GetTasksHandler).Methods("GET")
	router.HandleFunc("/tasks/{id}", handlers.GetTaskByIDHandler).Methods("GET")
	router.HandleFunc("/tasks/{id}", handlers.DeleteTaskHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}
