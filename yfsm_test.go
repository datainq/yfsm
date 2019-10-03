package yfsm

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	db       *sql.DB
	database = flag.Bool("database", false, "use database in tests")
)

func setupOnce() {
	c := "database=statemachinetest port=5433 user=statemachine password=statemachine sslmode=disable"
	var err error
	db, err = sql.Open("postgres", c)
	if err != nil {
		log.Fatal("cannot connect to Postgre")
	}
}

func tearDownOnce() {
	_ = db.Close()
}

func TestMain(m *testing.M) {
	flag.Parse()
	if *database {
		setupOnce()
	}
	code := m.Run()
	if *database {
		tearDownOnce()
	}
	os.Exit(code)
}

func TestSqlType_Get(t *testing.T) {
	if !*database {
		t.Log("skip database")
		return
	}

	tp := NewSqlType(db, "orders", "state_id")
	stateID, err := tp.Get(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, stateID)
}

func TestSqlType_Transition(t *testing.T) {
	if !*database {
		return
	}

	tp := NewSqlType(db, "orders", "state_id")
	assert.Error(t, ErrCannotFindInstance, tp.Transition(1, 2, 3))
	assert.NoError(t, tp.Transition(1, 1, 2))
	assert.NoError(t, tp.Transition(1, 2, 1))
}

func TestMachine_Can(t *testing.T) {
	if !*database {
		t.Log("skip database")
		return
	}
	tp := NewMapType()
	for i := 1; i < 11; i++ {
		assert.NoError(t, tp.Add(i, i))
	}
	machine := NewMachine(db, tp)
	for _, v := range []struct {
		id, from, event int
		ok              bool
	}{
		{1, 1, 1, true},
		{2, 2, 1, true},
		{3, 3, 1, true},
		{4, 4, 1, true},
		{5, 5, 1, true},
		{6, 6, 1, true},
		{7, 7, 1, true},
		{8, 8, 1, true},
		{9, 9, 1, true},
		{10, 10, 1, false},
	} {
		ok, err := machine.Can(v.id, EventFromID(v.event))
		assert.NoError(t, err)
		assert.Equal(t, v.ok, ok)

		toState, err := machine.ToState(v.id, EventFromID(1))
		if v.ok {
			assert.NoError(t, err)
			assert.Equal(t, v.from+1, toState)
		} else {
			assert.Error(t, err)
			assert.Equal(t, 0, toState)
		}
	}
}
