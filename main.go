package main

import (
	"BotWorks/src/Daemon"
	"BotWorks/src/Helpers"
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
	err := Helpers.SetWindowTitle("BotBuddy, v" + Helpers.Version)
	if err != nil {
		// silent
	}

	r := mux.NewRouter()

	// pages with templates
	r.HandleFunc("/", Handlers.IndexHandler)
	r.HandleFunc("/dashboard", Handlers.DashboardHandler)
	r.HandleFunc("/accounts", Handlers.AccountsHandler)
	r.HandleFunc("/proxies", Handlers.ProxyHandler)
	r.HandleFunc("/groups", Handlers.GroupsHandler)
	r.HandleFunc("/schedules", Handlers.SchedulesHandler)
	r.HandleFunc("/system", Handlers.SystemConfigHandler)
	r.HandleFunc("/scripts", Handlers.ScriptsHandler)
	r.HandleFunc("/massactions", Handlers.MassActionsHandler)
	r.HandleFunc("/logs", Handlers.LogsHandler)

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
		err := http.ListenAndServe(":80", r)
		if err != nil {
			println(err)
		}
	}()

	go Daemon.Begin()

	println("  ____        _     ____            _     _       ")
	println(" |  _ \\      | |   |  _ \\          | |   | |      ")
	println(" | |_) | ___ | |_  | |_) |_   _  __| | __| |_   _ ")
	println(" |  _ < / _ \\| __| |  _ <| | | |/ _` |/ _` | | | |")
	println(" | |_) | (_) | |_  | |_) | |_| | (_| | (_| | |_| |")
	println(" |____/ \\___/ \\__| |____/ \\__,_|\\__,_|\\__,_|\\__, |")
	println("                                             __/ |")
	println("   Version " + Helpers.Version + ", by lolwat                  |___/ ")
	println()
	println()

	Storage.RunUpdates()

	println("Server is now listening on port 80")

	// Keep this bitch going forever and ever, perhaps implement commands
	for true {
		text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		println(text)
	}
}
