package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/datainq/yfsm"
)

func main() {
	db, err := sql.Open("postgres", "database=statemachinetest port=5433 user=statemachine password=statemachine sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to Postgre")
	}

	// Assuming we have table orders with fields: id int, state int.
	t := yfsm.NewSqlType(db, "orders", "state_id")
	m := yfsm.NewMachine(db, t)

	created := yfsm.EventFromName("OK")
	ok, err := m.Can(1, created)
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		if err = m.Fire(17, created); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Print("cannot change state of object")
	}
}
