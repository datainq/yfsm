# yfsm

Yet (another) Finite State Machine library.

The library was created for one the ecommerce projects in DataInq ecommerce lab.

Main concepts are:
 - state
 - event

The graph of states and transitions can be described by tuples:
`<state0, event, state1>`.

The library has the capability of saving the state change history.
You should provide a `History` implementation.

## Usage

1. Create necessary tables in Postgre database.
2. Populate tables with states and data. One may use a tool creating
tables from tuples, check: `cmd/populate/populate.go`.
3. In your code, create a `Machine` instance with a `Type` object. You can use
   `SqlType` with table name and state column provided.
4. You can play with the states:

```
m := yfsm.NewMachine()
```