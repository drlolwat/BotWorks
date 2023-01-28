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
	r.HandleFunc("/start/id/{id}", Handlers.HandleStartBotById)
	r.HandleFunc("/start/group/{id}", Handlers.HandleStartBotsByGroup)
	r.HandleFunc("/stop/id/{id}", Handlers.HandleStopBotById)
	r.HandleFunc("/stop/group/{id}", Handlers.HandleStopBotByGroup)
	r.HandleFunc("/start/all", Handlers.HandleStartAll)
	r.HandleFunc("/stop/all", Handlers.HandleStopAll)

	fs := http.FileServer(http.Dir("assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	// begin routine for web serving
	go func() {
		err := http.ListenAndServe(":8080", r)
		if err != nil {
			println(err)
		}
	}()

	println("BotWorks is now running on port 8080.")

	// Keep this bitch going forever and ever, perhaps implement commands
	for true {
		text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		println(text)
	}
}
