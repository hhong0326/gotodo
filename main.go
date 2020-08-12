package main

import (
	"log"
	"net/http"

	"web7_todoweb/app"
)

func main() { // package보다 main에 있는게 좋다

	m := app.MakeHandler("./test.db") // flag. 실행인자로 가져오기
	defer m.Close()                   // model에 있는 db close 제어

	log.Println("Started App")
	err := http.ListenAndServe(":3000", m)

	if err != nil {
		panic(err)
	}
}
