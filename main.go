package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"github.com/niktheblak/passwordhash/pkg/scrypt"
	"html/template"
	"net/http"
	"os"
)

var (
	db    *sql.DB
	store *sessions.CookieStore
	tpl   *template.Template
)

const (
	sessionName = "session-name"
	userIDKey   = "user-id"
)

// User represents a user in the system
type User struct {
	ID       int
	Handle   string
	Password string
}

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./Botworks.db")
	if err != nil {
		println(err)
	}

	store = sessions.NewCookieStore([]byte("session"))
	tpl = template.Must(template.ParseGlob("./templates/*.html"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/dashboard", dashboardHandler)
	r.HandleFunc("/logout", logoutHandler)

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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}

	if isSessionValid(r) {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	}

	if r.PostForm.Has("handle") && r.PostForm.Has("password") {
		username := r.FormValue("handle")
		password := r.FormValue("password")

		if len(username) > 0 && len(password) > 0 {
			var user User
			if err := db.QueryRow("SELECT id, handle, password FROM customers WHERE handle = ?", username).Scan(&user.ID, &user.Handle, &user.Password); err != nil {
				println("login failed for user: " + username)
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			// Compare the provided password with the hashed password from the database
			if err := scrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
				println("login failed for user: " + username)
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			session, _ := store.Get(r, sessionName)
			session.Values[userIDKey] = user.ID
			session.Save(r, w)

			println("login was a success for user: " + username)
			http.Redirect(w, r, "/dashboard", http.StatusFound)
			return
		}
	}

	err = tpl.ExecuteTemplate(w, "bootstrap_header.html", nil)
	if err != nil {
		println("failed to execute tpl: bootstrap_header")
	}

	err = tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		println("failed to execute tpl: index")
	}
}

type DashboardData struct {
	RunningBots   int
	StoppedBots   int
	UniqueScripts int
	CurrentTime   string
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if !isSessionValid(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var user User
	session, _ := store.Get(r, sessionName)
	userID := session.Values[userIDKey]
	if userID != nil {
		user.ID = userID.(int)
	} else {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}

	fmt.Println("userid: ", user.ID)

	// get bot data for dash
	rows, err := db.Query("SELECT name,proxy_id,script_id,pid,group_id FROM accounts WHERE customer_id = ?", user.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer rows.Close()

	var runningBots = 0
	var stoppedBots = 0
	var uniqueScripts = 0
	var scripts []int

	for rows.Next() {
		var name string
		var proxyId int
		var scriptId int
		var pid int
		var groupId int

		err = rows.Scan(&name, &proxyId, &scriptId, &pid, &groupId)
		if err != nil {
			fmt.Println(err)
			return
		}

		if !containsInt(scriptId, scripts) {
			scripts = append(scripts, scriptId)
		}

		if pid > 0 {
			runningBots += 1
		} else {
			stoppedBots += 1
		}
	}

	uniqueScripts = len(scripts)

	fmt.Println("Bots running:", runningBots)
	fmt.Println("Bots offline:", stoppedBots)

	err = tpl.ExecuteTemplate(w, "bootstrap_header.html", nil)
	if err != nil {
		println("failed to execute tpl: bootstrap_header")
	}

	data := DashboardData{runningBots, stoppedBots, uniqueScripts, "1234"}
	err = tpl.ExecuteTemplate(w, "dashboard.html", data)
	if err != nil {
		println("failed to execute tpl: dashboard")
		return
	}

}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
	return
}

func isSessionValid(r *http.Request) bool {
	session, _ := store.Get(r, sessionName)
	userID := session.Values[userIDKey]
	if userID == nil {
		return false
	}

	return true
}

func containsInt(check int, numbers []int) bool {
	for _, item := range numbers {
		if item == check {
			return true
		}
	}
	return false
}
