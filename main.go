package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type person struct {
	id       int
	name     string
	position string
}

func main() {
	db, err := sql.Open("postgres", "user=postgres password=Utana_08 dbname=List sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	rows, err := db.Query("select * from employees")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	persons := []person{}
	for rows.Next() {
		p := person{}
		err := rows.Scan(&p.id, &p.name, &p.position)
		if err != nil {
			fmt.Println(err)
			continue
		}
		persons = append(persons, p)
	}

	for _, p := range persons {
		fmt.Println(p.id, p.name, p.position)

		_, err := db.Exec(
			"insert into employees2 (id, name, position) values ($1, $2, $3)",
			p.id,
			p.name,
			p.position,
		)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
