package Sessions

import "github.com/gorilla/sessions"

var (
	Store *sessions.CookieStore
)

const (
	SessionName = "session-name"
	UserIDKey   = "user-id"
)

func Init() {
	Store = sessions.NewCookieStore([]byte("session"))
}
