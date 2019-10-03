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
2. Populate tables with states and data.
3. In your code, create a `Machine` instance with a `Type` object. You can use
   `SqlType` with table name and state column provided.
4. You can play with the states:

```
m := yfsm.NewMachine()
```

## Configuring state machine

`state_machine` table contains all state machines you want to use. If you have
many object types you want to manage state, e.g. orders and warehouse
commodity transfer, you probably should create two separate machines:
`order_state` and `parcel_state`.

```
INSERT INTO state_machine(id, name) VALUES (1, 'order_states');
```

You must provide states:
```
INSERT INTO state_machine_state(id, state_machine_id, name, start, stop) VALUES
(1, 1, 'CREATED', TRUE, FALSE),
(2, 1, 'BUYER_DATA', FALSE, FALSE),
(3, 1, 'PAYMENT', FALSE, FALSE),
(4, 1, 'EXTERNAL_PAYMENT', FALSE, FALSE),
(5, 1, 'WAITING_FOR_ACCEPT', FALSE, FALSE),
(6, 1, 'COLLECT', FALSE, FALSE),
(7, 1, 'SEND', FALSE, FALSE),
(8, 1, 'ON_THE_WAY', FALSE, FALSE),
(9, 1, 'DELIVERED', FALSE, FALSE),
(10, 1, 'DONE', FALSE, TRUE),
(11, 1, 'PAYMENT_FAIL', FALSE, TRUE),
(12, 1, 'FAIL_TO_DELIVER', FALSE, TRUE),
(13, 1, 'DECLINED', FALSE, TRUE),
(14, 1, 'CANCELED', FALSE, TRUE);
```

Having states (think about them as nodes in graph), one must define transitions
between states. We have two concepts around that: `event` and `transition`.
The first can be though as a name e.g. `OK` may be an event, whereas transition
is a tuple of `(fromState, event, toState)`. There may be multiple transitions
defined for same event.

Let's define some events:
```
INSERT INTO state_machine_event(id, state_machine_id, name)
VALUES
(1, 1, 'OK'),
(2, 1, 'FAIL'),
(3, 1, 'RETRY');
```

Then we may define exact transitions we want to allow:
```
INSERT INTO state_machine_transition(
    id,
    state_machine_id,
    state_machine_event_id,
    from_state_id,
    to_state_id
) VALUES
(1, 1, 1, 1, 2),  -- 'CREATED' - on 'OK' change to 'BUYER_DATA'
(12, 1, 3, 11, 4);  -- 'DECLINED' - on 'RETRY' change to 'EXTERNAL_PAYMENT'
```
