CREATE TABLE state_machine
(
    id         SERIAL PRIMARY KEY NOT NULL,
    name       VARCHAR            NOT NULL UNIQUE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE state_machine_state
(
    id               SERIAL PRIMARY KEY                NOT NULL,
    state_machine_id INT REFERENCES state_machine (id) NOT NULL,
    name             VARCHAR                           NOT NULL,
    start            BOOL                              NOT NULL, -- potential start state,
    stop             BOOL                              NOT NULL, -- potential end state,

    created_at       TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at       TIMESTAMP WITH TIME ZONE
);

ALTER TABLE state_machine_state
    ADD CONSTRAINT state_machine_state_uniq_name_in_machine UNIQUE (state_machine_id, name);

CREATE TABLE state_machine_event
(
    id               SERIAL PRIMARY KEY                NOT NULL,
    state_machine_id INT REFERENCES state_machine (id) NOT NULL,
    name             VARCHAR                           NOT NULL,

    created_at       TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at       TIMESTAMP WITH TIME ZONE
);

ALTER TABLE state_machine_event
    ADD CONSTRAINT state_machine_event_uniq_name_in_machine UNIQUE (state_machine_id, name);

CREATE TABLE state_machine_transition
(
    id                     SERIAL PRIMARY KEY                      NOT NULL,
    state_machine_id       INT REFERENCES state_machine (id)       NOT NULL,
    state_machine_event_id INT REFERENCES state_machine_event (id) NOT NULL,
    from_state_id         INT REFERENCES state_machine_state (id) NOT NULL,
    to_state_id           INT REFERENCES state_machine_state (id) NOT NULL,

    created_at             TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at             TIMESTAMP WITH TIME ZONE
);

ALTER TABLE state_machine_transition
    ADD CONSTRAINT state_machine_transition_uniq_ UNIQUE (state_machine_id, state_machine_event_id, from_state_id);
