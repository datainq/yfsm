package main

import (
	"database/sql"
	"log"

	"github.com/datainq/yfsm"
)

func main() {
	var db *sql.DB
	// Assuming we have table orders with fields: id int, state int.
	t := yfsm.NewSqlType(db, "orders", "state")
	m := yfsm.NewMachine(db, t)

	created := yfsm.EventFromName("CREATED")
	ok, err := m.Can(17, created)
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
