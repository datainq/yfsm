CREATE TABLE state_machine
(
    id         SERIAL PRIMARY KEY NOT NULL,
    name       VARCHAR            NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE state_machine_state
(
    id         SERIAL PRIMARY KEY NOT NULL,
    model      VARCHAR            NOT NULL,
    name       VARCHAR            NOT NULL,
    start      BOOL               NOT NULL, -- potential start state,
    stop       BOOL               NOT NULL, -- potential end state,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE state_machine_event
(
    id               SERIAL PRIMARY KEY               NOT NULL,
    state_machine_id INT REFERENCES state_machine (id) NOT NULL,

    created_at       TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at       TIMESTAMP WITH TIME ZONE
);

CREATE TABLE state_machine_transition
(
    id                  SERIAL PRIMARY KEY                    NOT NULL,
    state_machine_id    INT REFERENCES state_machine (id)      NOT NULL,
    state_machine_event_id INT REFERENCES state_machine_event (id) NOT NULL,
    start_state_id      INT REFERENCES state_machine_state (id) NOT NULL,
    end_state_id        INT REFERENCES state_machine_state (id) NOT NULL,

    created_at          TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at          TIMESTAMP WITH TIME ZONE
);
