package yfsm

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gocraft/dbr/v2"
	"github.com/gocraft/dbr/v2/dialect"
	"github.com/lib/pq"
)

type History interface {
	Save()
}

type Storage interface {
	Get()
	Save()
}

type Machine interface {
	ID() int
	Name() string

	// Check if a type's object is in a proper state to handle event.
	Can(id int, event Event) (bool, error)
	Fire(id int, event Event) error
}

type State interface {
	ID() int
	Machine() int
	Name() string
}

// Event describes an event to happen.
type Event interface {
	ID() int
	Transition() int
	Name() string
	// Machine() int
	// FromState() int
	// ToState() int

	Identify() bool
}

func EventFromName(name string) Event {
	return event{name: name}
}

func EventFromTransition(id int) Event {
	return event{transition: id}
}

func EventFromID(id int) Event {
	return event{id: id}
}

type event struct {
	id         int
	transition int
	machine    int
	name       string
	fromState  int
	toState    int
}

func (e event) ID() int {
	return e.id
}

func (e event) Transition() int {
	return e.transition
}

func (e event) Machine() int {
	return e.machine
}

func (e event) Name() string {
	return e.name
}

func (e event) FromState() int {
	return e.fromState
}

func (e event) ToState() int {
	return e.toState
}

// Identify checks if a provided data identifies transition.
// It must be one of:
//  - the state_machine_transition_id
//  - from_state, state_machine_event_id
func (e event) Identify() bool {
	return e.transition > 0 || e.id > 0 && e.fromState > 0
}

type SqlType struct {
	db               *sql.DB
	selectQ, updateQ string
}

func NewSqlType(db *sql.DB, table string, column string) *SqlType {
	selectQ := fmt.Sprintf("SELECT %s FROM %s WHERE id=$1",
		pq.QuoteIdentifier(column),
		pq.QuoteIdentifier(table))
	updateQ := fmt.Sprintf("UPDATE %s SET %s=$1 WHERE id=$2 AND %s=$3 RETURNING id",
		pq.QuoteIdentifier(table), pq.QuoteIdentifier(column),
		pq.QuoteIdentifier(column))
	return &SqlType{db: db, selectQ: selectQ, updateQ: updateQ}
}

func (s *SqlType) Get(id int) (int, error) {
	err := s.db.QueryRow(s.selectQ, id).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, ErrCannotFindInstance
	}
	return id, err
}

func (s *SqlType) Transition(id, fromState, toState int) error {
	err := s.db.QueryRow(s.updateQ, toState, id, fromState).Scan(&id)
	if err == sql.ErrNoRows {
		return ErrCannotFindInstance
	}
	return err
}

type Type interface {
	Get(id int) (int, error)
	Transition(id, fromState, toState int) error
}

// Loads a given machine for a given type.
func LoadMachineForType(id int, t Type) Machine {
	return nil
}

type machine struct {
	id   int
	name string
	t    Type
	db   *dbr.Session
}

func (m machine) ID() int {
	panic("implement me")
}

func (m machine) Name() string {
	panic("implement me")
}

var (
	ErrCannotFindInstance   = errors.New("cannot find object")
	ErrCannotFindTransition = errors.New("cannot find transition")
	ErrCannotIdentifyEvent  = errors.New("cannot identify event")
)

// Can verifies that event can be handled. It queries database for a transition
// where state matches the resource state and state_machine_id matches
// the machine. Then:
//  1. if transition id is specified, we query for it.
//  2. if event id is specified, we check for event id.
//  3. if event name is specified, we look for transition under such name.
func (m machine) Can(id int, event Event) (bool, error) {
	stateID, err := m.t.Get(id)
	if err != nil {
		return false, err
	}
	_, err = m.ToState(stateID, event)
	return err == nil, err
}

func (m machine) ToState(fromState int, event Event) (int, error) {
	q := m.db.Select("end_state_id").
		From("state_machine_transition").
		Where(dbr.Eq("from_state_id", fromState),
			dbr.Eq("state_machine_id", m.id))

	var hasCheck bool
	if transition := event.Transition(); transition > 0 {
		q.Where("id=?", transition)
		hasCheck = true
	}
	if id := event.ID(); id > 0 {
		q.Where("state_machine_event_id=?", event.ID())
		hasCheck = true
	}
	if name := event.Name(); name != "" {
		q.Where(dbr.Expr("state_machine_event_id IN (SELECT id FROM state_machine_event WHERE name=?)", name))
		hasCheck = true
	}
	if !hasCheck {
		return 0, ErrCannotIdentifyEvent
	}

	var endStateID int
	err := q.LoadOne(&endStateID)
	if err == dbr.ErrNotFound {
		return 0, ErrCannotFindTransition
	}
	return endStateID, err
}

func (m machine) Fire(id int, event Event) error {
	fromState, err := m.t.Get(id)
	if err != nil {
		return err
	}
	toState, err := m.ToState(fromState, event)
	if err != nil {
		return err
	}
	return m.t.Transition(id, fromState, toState)
}

func NewMachine(rawDB *sql.DB, t Type) Machine {
	db := &dbr.Connection{
		DB:            nil,
		Dialect:       dialect.PostgreSQL,
		EventReceiver: &dbr.NullEventReceiver{},
	}
	s := db.NewSession(nil)
	return &machine{db: s, t: t}
}
