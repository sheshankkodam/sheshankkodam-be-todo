package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"time"
)

type Database struct {
	psql *sql.DB
}

func Open(dbUrl string) (*Database, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	return &Database{psql:db}, nil
}

func (d *Database) Init() error {
	if err := d.createTables(); err != nil {
		return err
	}
	return nil
}

func (d *Database) createTables() error {
	if _, err := d.psql.Exec("CREATE TABLE IF NOT EXISTS customer (customer_id VARCHAR(255) PRIMARY KEY, username VARCHAR(255), password VARCHAR(255), last_login TIMESTAMP)"); err != nil {
		return fmt.Errorf("error creating customer table, err=%s", err)
	}

	if _, err := d.psql.Exec("CREATE TABLE IF NOT EXISTS todo (todo_id UUID NOT NULL, customer_id VARCHAR(255), name VARCHAR(255), completed BOOLEAN, priority VARCHAR(10), created_at TIMESTAMP, PRIMARY KEY (todo_id, customer_id), FOREIGN KEY (customer_id) REFERENCES customer (customer_id))"); err != nil {
		return fmt.Errorf("error creating customer table, err=%s", err)
	}
	return nil
}

func (d *Database) InsertCustomer(customerId, username, pwd string) error {
	insertQuery := `INSERT INTO customer (customer_id, username, password, last_login)VALUES ($1, $2, $3, $4)`
	if _, err := d.psql.Exec(insertQuery, customerId, username, pwd, time.Now()); err != nil {
		return fmt.Errorf("error inserting customer, err=%s", err)
	}
	return nil
}

func (d *Database) GetTodoTasks(customerId string) (*sql.Rows, error) {
	selectQuery := `SELECT * FROM todo WHERE customer_id=$1`
	rows, err := d.psql.Query(selectQuery, customerId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving data for customer=%s err=%s", customerId, err)
	}
	return rows, nil
}

func (d *Database) InsertTodo(todoId uuid.UUID, customerId, name, priority string, completed bool) (time.Time, error) {
	insertQuery := `INSERT INTO todo (todo_id, customer_id, name, completed, priority, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	timeNow := time.Now()
	if _, err := d.psql.Exec(insertQuery, todoId, customerId, name, completed, priority, timeNow); err != nil {
		return timeNow, fmt.Errorf("error inserting todo, err=%s", err)
	}
	return timeNow, nil
}

func (d *Database) DeleteTodo(todoId uuid.UUID) (sql.Result, error) {
	delQuery := `DELETE FROM todo WHERE todo_id=$1`
	res, err := d.psql.Exec(delQuery, todoId)
	if err != nil {
		return nil, fmt.Errorf("error deleting todo, todo=%s err=%s", todoId, err)
	}
	return res, nil
}

func (d *Database) UpdateTodo(todoId uuid.UUID, customerId string, name, priority string, completed bool) (time.Time, error) {

	updQuery := `UPDATE todo SET name = $3, completed = $4, priority = $5, created_at = $6 WHERE todo_id = $1 AND customer_id = $2;`
	timeNow := time.Now()

	_, err := d.psql.Exec(updQuery, todoId, customerId, name, completed, priority, timeNow)
	if err != nil {
		return timeNow, fmt.Errorf("error updating todo=%s customer=%s, err=%s", todoId, customerId, err)
	}
	return timeNow, nil
}
