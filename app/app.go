package app

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"web7_todoweb/model"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

//암호화 되는 쿠키 저장
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var rd *render.Render

type AppHandler struct {
	http.Handler
	db model.DBHandler
}

//func 포인터를 갖는 variable
var getSessionID = func(r *http.Request) string {
	session, err := store.Get(r, "session")
	if err != nil {
		return ""
	}

	val := session.Values["id"]
	str := fmt.Sprintf("%v", val)

	if val == nil {
		return ""
	}

	return str
}

func (a *AppHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func (a *AppHandler) getToDoListHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionID(r)
	list := a.db.GetTodos(sessionID)

	// list := []*model.Todo{}
	// for _, v := range todoMap {
	// 	list = append(list, v)
	// }

	rd.JSON(w, http.StatusOK, list)
}

func (a *AppHandler) addToDoListHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionID(r)
	name := r.FormValue("name")
	todo := a.db.AddTodo(sessionID, name)
	// id := len(todoMap) + 1
	// todo := &Todo{id, name, false, time.Now()}
	// todoMap[id] = todo
	//의존성 분리하는 작업~ 메모리 쓰는 부분을 다른 패키지로 옮김!

	rd.JSON(w, http.StatusCreated, todo)
}

func (a *AppHandler) removeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	ok := a.db.RemoveTodo(id)

	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}

	//map 확인
	// if _, ok := todoMap[id]; ok {
	// 	delete(todoMap, id)
	// 	rd.JSON(w, http.StatusOK, Success{true})
	// } else {
	// 	rd.JSON(w, http.StatusOK, Success{false})
	// }
}

func (a *AppHandler) completeTodoHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	completed := r.FormValue("complete") == "true"

	ok := a.db.CompleteTodo(id, completed)

	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}

	// if todo, ok := todoMap[id]; ok {
	// 	todo.Comleted = completed

	// 	rd.JSON(w, http.StatusOK, Success{true})
	// } else {
	// 	rd.JSON(w, http.StatusOK, Success{false})
	// }

}

type Success struct {
	Success bool `json:"success"`
}

func (a *AppHandler) Close() {
	a.db.Close()
}

func CheckSignin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if strings.Contains(r.URL.Path, "/signin") ||
		strings.Contains(r.URL.Path, "/auth") {
		next(w, r)
		return
	}

	//if user already signed in
	sessionID := getSessionID(r)
	if sessionID != "" {
		next(w, r)
		return
	}

	//else
	//redirect signin.html
	http.Redirect(w, r, "/signin.html", http.StatusTemporaryRedirect)
}

func MakeHandler(filepath string) *AppHandler { // handler는 늘 포인터?

	rd = render.New()
	r := mux.NewRouter()

	// n := negroni.Classic()
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), negroni.HandlerFunc(CheckSignin), negroni.NewStatic(http.Dir("public")))
	n.UseHandler(r)

	a := &AppHandler{
		Handler: n,
		db:      model.NewDBHandler(filepath),
	}

	r.HandleFunc("/", a.indexHandler)
	r.HandleFunc("/todos", a.getToDoListHandler).Methods("GET")
	r.HandleFunc("/todos", a.addToDoListHandler).Methods("POST")
	r.HandleFunc("/todos/{id:[0-9]+}", a.removeTodoHandler).Methods("DELETE")
	r.HandleFunc("/complete-todo/{id:[0-9]+}", a.completeTodoHandler).Methods("GET")

	r.HandleFunc("/auth/google/login", googleLoginHandler)
	r.HandleFunc("/auth/google/callback", googleAuthCallback)

	return a
}
