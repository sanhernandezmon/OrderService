package main

import (
	"OrderService/domain"
	"OrderService/repository"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/api/v1/order", createNewOrder).Methods("POST")
	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second argument
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func createNewOrder(w http.ResponseWriter, r *http.Request) {
	var request domain.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	orderID, err := repository.SaveOrderToDynamoDB(request)
	if err != nil {
		log.Fatal("error writting on dynamo db")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	repository.SendOrderSQSMessage(uuid.New().String(), request.TotalPrice)
	fmt.Fprintf(w, "{\"added\":\"%s\"}", orderID)
}
func main() {
	println("Order Service started")
	handleRequests()
}
