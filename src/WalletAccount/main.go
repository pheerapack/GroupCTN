package main

import (
	"log"
	"net/http"
	"WalletAccount/Account"
	"gopkg.in/gorilla/handlers"
	"gopkg.in/mux"
	"time"
	"fmt"
	"gopkg.in/mgo.v2"
)

func main() {
	/*
	router := Account.NewRouter() // create routes
	// these two lines are important in order to allow access from the front-end side to the methods
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	// launch server with CORS validations
	log.Fatal(http.ListenAndServe(":5000",	handlers.CORS(allowedOrigins, allowedMethods)(router)))
	*/
	router := mux.NewRouter()
	//people = append(people, Person{ID: "1", Firstname: "Nic", Lastname: "Raboy", Address: &Address{City: "Dublin", State: "CA"}})
	//people = append(people, Person{ID: "2", Firstname: "Maria", Lastname: "Raboy"})
	//router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	//router.HandleFunc("/people/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/WalletAccount/{id}", CreateWalletAccount).Methods("POST")
	//router.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func CreateWalletAccount() {
	time1 := time.Now()
	fmt.Println(time1)
	session, err := mgo.Dial("0.0.0.0:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	datasource := New()
	buf := make([]byte, 32)
	name_, _ := datasource.Read(buf)
	phone_, _ := datasource.Read(buf)
	name := string(name_)
	phone := string(phone_)
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("people")
	//people := []*Person{}
	for i := 0; i < 1000000; i++ {
		c.Insert(&Person{name, phone})
		if i%10000 == 0 {
			log.Printf("Populated %d records", i)
		}
	}
	fmt.Printf("done")
	time2 := time.Since(time1)
	fmt.Println(time2)
