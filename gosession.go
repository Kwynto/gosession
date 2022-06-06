package gosession

import (
	"crypto/rand"
	"net/http"
)

const (
	GOSESSION_COOKIE_NAME string = "SessionId"
	GOSESSION_MAX_AGE     int    = 43_200 // Max age is 12 hours.
)

type SessionId string
type Session map[string]interface{}
type Sessions map[SessionId]Session

// TODO: Сделать очистку сервеного хранилища сессий от старых записей
var AllSessions Sessions

// Privat

func generateId() SessionId {
	b := make([]byte, 32)
	rand.Read(b)
	return SessionId(b)
}

func getOrSetCookie(w *http.ResponseWriter, r *http.Request) SessionId {
	data, err := r.Cookie(GOSESSION_COOKIE_NAME)
	if err != nil {
		id := generateId()
		cookie := &http.Cookie{
			Name:   GOSESSION_COOKIE_NAME,
			Value:  string(id),
			MaxAge: 0,
		}
		http.SetCookie(*w, cookie)
		return id
	}
	return SessionId(data.Value)
}

func deleteCookie(w *http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   GOSESSION_COOKIE_NAME,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(*w, cookie)
}

// Public

func (id SessionId) Set(name string, value interface{}) {
	ses := AllSessions[id]
	ses[name] = value
	AllSessions[id] = ses
}

func (id SessionId) GetAll() Session {
	return AllSessions[id]
}

func (id SessionId) GetOne(name string) interface{} {
	data := AllSessions[id]
	return data[name]
}

func (id SessionId) RemoveSession(w *http.ResponseWriter) {
	delete(AllSessions, id)
	deleteCookie(w)
}

func (id SessionId) RemoveValue(name string) {
	data := AllSessions[id]
	delete(data, name)
	AllSessions[id] = data
}

func Start(w *http.ResponseWriter, r *http.Request) SessionId {
	id := getOrSetCookie(w, r)
	data := AllSessions[id]
	if data == nil {
		data := make(Session, 0)
		AllSessions[id] = data
	}
	return id
}

func init() {
	AllSessions = make(Sessions, 0)
}
