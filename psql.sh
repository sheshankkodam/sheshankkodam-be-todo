#!/usr/bin/env bash

CREATE DATABASE todo;
CREATE TABLE IF NOT EXISTS customer (customer_id VARCHAR(255) PRIMARY KEY, username VARCHAR(255), password VARCHAR(255), last_login TIMESTAMP);

INSERT INTO customer (customer_id, username, password, last_login) VALUES ("123", "sheshank". "encryptedpwd", timstamp);


CREATE TABLE task (todo_id INTEGER NOT NULL, customer_id VARCHAR(255), name VARCHAR(255), priority VARCHAR(10), created_at TIMESTAMP, PRIMARY KEY (todo_id, customer_id), FOREIGN KEY (customer_id) REFERENCES customer (customer_id))
