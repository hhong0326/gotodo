package model

import (
	"time"
)

type Todo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type DBHandler interface {
	GetTodos(sessionId string) []*Todo
	AddTodo(sessionId string, name string) *Todo
	RemoveTodo(id int) bool
	CompleteTodo(id int, completed bool) bool
	Close()
}

func NewDBHandler(filepath string) DBHandler {
	// handler = newMemoryHandler()
	return newSqliteHandler(filepath)
}

// func init() {
// 	// handler = newMemoryHandler()
// 	handler = newSqliteHandler()
// }

// func GetTodos() []*Todo {

// 	return handler.getTodos()
// }
// func AddTodo(name string) *Todo {

// 	return handler.addTodo(name)
// }

// func RemoveTodo(id int) bool {

// 	return handler.removeTodo(id)
// }

// func CompleteTodo(id int, completed bool) bool {

// 	return handler.completeTodo(id, completed)
// }
