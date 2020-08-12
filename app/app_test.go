package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/hhong0326/gotodo/model"

	"github.com/stretchr/testify/assert"
)

func TestTodos(t *testing.T) {

	//func 포인터를 바꿔버려서 소스 수정없이 테스트를 가능하게 한다
	getSessionID = func(r *http.Request) string {
		return "testsessionId"
	}
	os.Remove("./test.db")
	assert := assert.New(t)
	ah := MakeHandler("./test.db")
	defer ah.Close()

	ts := httptest.NewServer(ah)
	defer ts.Close()

	res, err := http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test To do"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)

	var todo model.Todo
	err = json.NewDecoder(res.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test To do")
	id1 := todo.ID

	res, err = http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test To do2"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)

	var todo2 model.Todo
	err = json.NewDecoder(res.Body).Decode(&todo2)
	assert.NoError(err)
	assert.Equal(todo2.Name, "Test To do2")
	id2 := todo2.ID

	res, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	todos := []*model.Todo{}
	json.NewDecoder(res.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)

	// index, value := range []
	for _, t := range todos {
		if t.ID == id1 {
			assert.Equal("Test To do", t.Name)
		} else if t.ID == id2 {
			assert.Equal("Test To do2", t.Name)
		} else {
			assert.Error(fmt.Errorf("testID should be id1 or id2"))
		}
	}

	res, err = http.Get(ts.URL + "/complete-todo/" + strconv.Itoa(id1) + "?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	//
	res, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	todos = []*model.Todo{}
	json.NewDecoder(res.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)

	// index, value := range []
	for _, t := range todos {
		if t.ID == id1 {
			assert.True(t.Completed)
		}
	}

	//delete는 만들어야함
	req, _ := http.NewRequest("DELETE", ts.URL+"/todos/"+strconv.Itoa(id1), nil)

	res, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	//list
	res, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	// var todos []*Todo
	json.NewDecoder(res.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 1)

	// index, value := range []
	for _, t := range todos {
		if t.ID == id1 {
			assert.Equal(t.ID, id2)
		}
	}

}
