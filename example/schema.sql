-- create database and user
CREATE DATABASE statemachinetest;
CREATE USER statemachine WITH ENCRYPTED PASSWORD 'statemachine';
GRANT ALL PRIVILEGES ON DATABASE statemachinetest TO statemachine;

-- create test table for orders
CREATE TABLE IF NOT EXISTS orders(
    id SERIAL PRIMARY KEY,
    state_id INT NOT NULL REFERENCES state_machine_state(id)
);

INSERT INTO state_machine(id, name) VALUES (1, 'order_states');

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

INSERT INTO state_machine_event(id, state_machine_id, name)
VALUES
(1, 1, 'OK'),
(2, 1, 'FAIL'),
(3, 1, 'RETRY');

INSERT INTO state_machine_transition(id, state_machine_id, state_machine_event_id, from_state_id, to_state_id)
VALUES
-- OK
(1, 1, 1, 1, 2),
(2, 1, 1, 2, 3),
(3, 1, 1, 3, 4),
(4, 1, 1, 4, 5),
(5, 1, 1, 5, 6),
(6, 1, 1, 6, 7),
(7, 1, 1, 7, 8),
(8, 1, 1, 8, 9),
(9, 1, 1, 9, 10),
(10, 1, 2, 4, 11),
(11, 1, 2, 8, 12),
(12, 1, 3, 11, 4),
(13, 1, 3, 12, 8),
(14, 1, 2, 5, 13);

INSERT INTO orders(id, state_id)
VALUES
 (1, 1),
 (2, 2),
 (3, 3),
 (4, 4),
 (5, 5),
 (6, 6),
 (7, 7),
 (8, 8),
 (9, 9),
 (10, 10),
 (11, 11),
 (12, 12),
 (13, 13),
 (14, 14);
