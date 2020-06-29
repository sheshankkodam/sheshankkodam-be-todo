package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	TodoId 	 	string 		`json:"todo_id"`
	CustomerId  string 		`json:"customer_id"`
	Name	 	string 		`json:"name"`
	Completed 	bool        `json:"completed"`
	Priority 	string 		`json:"priority"`
	CreatedAt   time.Time 	`json:"created_at"`
}

type DeleteToDoResp struct {
	TodoId 	 	string 	  `json:"todo_id"`
	CustomerId  string    `json:"customer_id"`
	RowsDeleted int64     `json:"rows_deleted"`
}

func (s *Server) addTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerId := vars["customerId"]

	var t Todo

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todoId := uuid.New()

	insertedTime, insErr := s.db.InsertTodo(todoId, customerId, t.Name, t.Priority, t.Completed)
	if insErr != nil {
		http.Error(w, insErr.Error(), http.StatusBadRequest)
		return
	}

	t.TodoId = todoId.String()
	t.CustomerId = customerId
	t.CreatedAt = insertedTime

	resp, mErr := json.Marshal(t)
	if mErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling task response.")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error writing task response.")
	}
	return
}

func (s *Server) getTodoTasks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerId := vars["customerId"]
	var todos []*Todo

	rows, err := s.db.GetTodoTasks(customerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	for rows.Next() {
		var todo_id uuid.UUID
		var customer_id string
		var name string
		var completed bool
		var priority string
		var created_at time.Time

		scanErr := rows.Scan(&todo_id, &customer_id, &name, &completed, &priority, &created_at)
		if scanErr != nil {
			http.Error(w, scanErr.Error(), http.StatusInternalServerError)
			return
		}

		todo := &Todo{
			TodoId:     todo_id.String(),
			CustomerId: customerId,
			Name:       name,
			Completed:     completed,
			Priority:   priority,
			CreatedAt:  time.Time{},
		}

		todos = append(todos, todo)
	}

	todosJson, err := json.Marshal(todos)
	if err != nil {
		log.Fatal("Cannot encode to JSON ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(todosJson); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error writing task response.")
	}
	return
}

func (s *Server) deleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerId, _ := vars["customerId"]
	todoId, _ := uuid.Parse(vars["todoId"])

	res, err := s.db.DeleteTodo(todoId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		http.Error(w, rowsErr.Error(), http.StatusInternalServerError)
		return
	}

	resp := &DeleteToDoResp{
		TodoId:      todoId.String(),
		CustomerId:  customerId,
		RowsDeleted: rowsAffected,
	}


	respJson, err := json.Marshal(resp)
	if err != nil {
		log.Fatal("Cannot encode to JSON ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(respJson); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error writing task response.")
	}
	return
}

func (s *Server) updateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerId, _ := vars["customerId"]
	todoId, _ := uuid.Parse(vars["todoId"])

	t := &Todo{
		TodoId:     todoId.String(),
		CustomerId: customerId,
	}

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedTime, err := s.db.UpdateTodo(todoId, customerId, t.Name, t.Priority, t.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.CreatedAt = updatedTime

	resp, mErr := json.Marshal(t)
	if mErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling task response.")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error writing task response.")
	}
	return
}

