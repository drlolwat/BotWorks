package main

import (
	"BotWorks/src/Storage"
	"BotWorks/src/Web/Handlers"
	"BotWorks/src/Web/Sessions"
	"BotWorks/src/Web/Template"
	"bufio"
	_ "fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"os"
)

func init() {
	Storage.Init()
	Sessions.Init()
	Template.Init()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Handlers.IndexHandler)
	r.HandleFunc("/dashboard", Handlers.DashboardHandler)
	r.HandleFunc("/logout", Handlers.LogoutHandler)

	fs := http.FileServer(http.Dir("assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	go func() {
		err := http.ListenAndServe(":8080", r)
		if err != nil {
			println(err)
		}
	}()

	println("BotWorks is now running on port 8080.")
	for true {
		text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		println(text)
	}
}
