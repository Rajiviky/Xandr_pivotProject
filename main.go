package main

import (
	"net/http"

	"pivot/pkg/github.com/gorilla/mux"
)

func main() {

	// route handeller
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/default", handlerFunc)
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("src/static/"))))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	router.HandleFunc("/home/", login)
	router.HandleFunc("/profile", profile)
	router.HandleFunc("/logout", logout)
	router.HandleFunc("/hello/", hello)

	// api call route
	router.HandleFunc("/api/employeeskill/{attuid}", GetEmployeeskill).Methods("GET")
	router.HandleFunc("/api/addemployee/", AddEmployee).Methods("POST")
	router.HandleFunc("/api/getemployee/", GetEmployee).Methods("GET")
	router.HandleFunc("/api/deleteemployee/{attuid}", DeleteEmployee).Methods("DELETE")
	router.HandleFunc("/api/addskill/", AddSkill).Methods("POST")
	router.HandleFunc("/api/deleteemployeeskill/{attuid-skillid}", DeleteEmployeeSkill).Methods("DELETE")

	// website url
	http.ListenAndServe("localhost:3000", router)
}
