package model

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3" // _ 암시적 사용
)

type sqliteHandler struct {
	db *sql.DB
}

func (s *sqliteHandler) GetTodos(sessionId string) []*Todo {
	todos := []*Todo{}
	rows, err := s.db.Query("SELECT id, name, completed, createdAt FROM todos WHERE sessionId=?", sessionId)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() { // Next가 true면 loop
		var todo Todo

		rows.Scan(&todo.ID, &todo.Name, &todo.Completed, &todo.CreatedAt)
		todos = append(todos, &todo)
	}

	return todos
}

func (s *sqliteHandler) AddTodo(sessionId string, name string) *Todo {
	stmt, err := s.db.Prepare("INSERT INTO todos (sessionId, name, completed, createdAt) VALUES (?, ?, ?, datetime('now'))")
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(sessionId, name, false)
	if err != nil {
		panic(err)
	}
	id, _ := rst.LastInsertId()
	var todo Todo
	todo.ID = int(id)
	todo.Name = name
	todo.Completed = false
	todo.CreatedAt = time.Now()
	return &todo
}

func (s *sqliteHandler) RemoveTodo(id int) bool {
	stmt, err := s.db.Prepare("DELETE FROM todos WHERE id=?")

	if err != nil {
		panic(err)
	}

	result, err := stmt.Exec(id)

	if err != nil {
		panic(err)
	}

	//쿼리로부터 영향받은 레코드 갯수
	cnt, _ := result.RowsAffected()

	return cnt > 0
}

func (s *sqliteHandler) CompleteTodo(id int, completed bool) bool {
	stmt, err := s.db.Prepare("UPDATE todos SET completed=? WHERE id=?")

	if err != nil {
		panic(err)
	}

	result, err := stmt.Exec(completed, id)

	if err != nil {
		panic(err)
	}

	//쿼리로부터 영향받은 레코드 갯수
	cnt, _ := result.RowsAffected()

	return cnt > 0
}

func (s *sqliteHandler) Close() {
	s.db.Close()
}

func newSqliteHandler(filepath string) DBHandler {
	database, err := sql.Open("sqlite3", filepath)

	if err != nil {
		panic(err)
	}

	statement, _ := database.Prepare(
		`CREATE TABLE IF NOT EXISTS todos (
			id        INTEGER  PRIMARY KEY AUTOINCREMENT,
			sessionId STRING,
			name      TEXT,
			completed BOOLEAN,
			createdAt DATETIME
		);
		CREATE INDEX IF NOT EXISTS sessionIdIndexOnTodos ON todos (
			sessionId ASC
		);`) // index를 만들어 많은 컬럼의 리딩 속도롤 log2N으로 만들기
	// where 절에 들어가는 필드값에는 index를 걸어주는게 seaching 속도 빠름!!
	statement.Exec()
	return &sqliteHandler{db: database}
}
