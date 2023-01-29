package main

import (
	"BotWorks/src/Storage"
	"BotWorks/src/Web/Handlers"
	"BotWorks/src/Web/Handlers/Actions"
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

	// pages with templates
	r.HandleFunc("/", Handlers.IndexHandler)
	r.HandleFunc("/dashboard", Handlers.DashboardHandler)
	r.HandleFunc("/accounts", Handlers.AccountsHandler)

	// pages that are essentially actions and redirects
	r.HandleFunc("/logout", Actions.LogoutAction)
	r.HandleFunc("/start/id/{id}", Actions.StartBotByIdAction)
	r.HandleFunc("/start/group/{id}", Actions.StartBotsByGroupAction)
	r.HandleFunc("/stop/id/{id}", Actions.StopBotByIdAction)
	r.HandleFunc("/stop/group/{id}", Actions.StopBotByGroupAction)
	r.HandleFunc("/start/all", Actions.StartAllBotsAction)
	r.HandleFunc("/stop/all", Actions.StopAllBotsAction)

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
