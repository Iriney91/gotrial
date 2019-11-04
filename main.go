package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Person struct {
	ID       int
	Name     string
	Position string
}

var (
	db         *sql.DB
	listStmt   *sql.Stmt
	singleStmt *sql.Stmt
)

func main() {
	var err error

	db, err = sql.Open("postgres", "user=postgres password=Utana_08 dbname=List sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	listStmt, err = db.Prepare("select * from employees")
	if err != nil {
		panic(err)
	}

	singleStmt, err = db.Prepare("select * from employees where id=$1")
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/employee", getEmployees).Methods("GET")
	r.HandleFunc("/employee/{id:[0-9]+}", getEmployeeByID).Methods("GET")

	srv := &http.Server{
		Addr:    "localhost:8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	srv.Shutdown(ctx)
	os.Exit(0)
}

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
