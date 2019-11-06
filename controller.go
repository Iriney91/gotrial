package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func getEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := listStmt.Query()
	if err != nil {
		log.Println(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	persons := []Person{}
	for rows.Next() {
		p := Person{}
		err := rows.Scan(&p.ID, &p.Name, &p.Position)
		if err != nil {
			log.Println(err)
			continue
		}
		persons = append(persons, p)
	}

	if err := json.NewEncoder(w).Encode(persons); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getEmployeeByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	time.Sleep(time.Second * 5)
	id := mux.Vars(r)["id"]
	p := Person{}
	if err := singleStmt.QueryRow(id).Scan(&p.ID, &p.Name, &p.Position); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
