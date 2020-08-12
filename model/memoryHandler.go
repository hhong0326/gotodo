package model

import "time"

type memoryHandler struct {
	todoMap map[int]*Todo
}

func (m *memoryHandler) GetTodos(sessionId string) []*Todo {
	list := []*Todo{}
	for _, v := range m.todoMap {
		list = append(list, v)
	}
	return list
}

func (m *memoryHandler) AddTodo(sessionId string, name string) *Todo {
	id := len(m.todoMap) + 1
	todo := &Todo{id, name, false, time.Now()}
	m.todoMap[id] = todo

	return todo
}

func (m *memoryHandler) RemoveTodo(id int) bool {
	if _, ok := m.todoMap[id]; ok {
		delete(m.todoMap, id)
		return true
	}

	return false
}

func (m *memoryHandler) CompleteTodo(id int, completed bool) bool {
	if todo, ok := m.todoMap[id]; ok {
		todo.Completed = completed
		return true
	}

	return false
}

func (m *memoryHandler) Close() {

}

func newMemoryHandler() DBHandler {
	m := &memoryHandler{} // memory핸들러가 db핸들러 interface를 implement 하고 있기 때문에
	m.todoMap = make(map[int]*Todo)

	return m
}
