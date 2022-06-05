package gosession

import (
	"crypto/rand"
	"net/http"
)

const COOKIE_NAME = "sessionId"

type sessionId string
type Session map[string]interface{}
type Sessions map[sessionId]Session

var AllSessions Sessions

// Privat

func generateId() sessionId {
	b := make([]byte, 32)
	rand.Read(b)
	// return fmt.Sprintf("%x", b)
	return sessionId(b)
}

// Public

func (id sessionId) GetAll() Session {
	return AllSessions[id]
}

func (id sessionId) GetOne(name string) interface{} {
	data := AllSessions[id]
	return data[name]
}

func (id sessionId) Set(name string, value interface{}) {
	ses := id.GetAll()
	ses[name] = value
	AllSessions[id] = ses
}

func getOrSetId(w *http.ResponseWriter, r *http.Request) sessionId {
	data, err := r.Cookie(COOKIE_NAME)
	if err != nil {
		gi := generateId()
		cookie := &http.Cookie{
			Name:   COOKIE_NAME,
			Value:  string(gi),
			MaxAge: 0,
		}
		http.SetCookie(*w, cookie)
		return gi
	}
	return sessionId(data.Value)
}

func Start(w *http.ResponseWriter, r *http.Request) sessionId {
	id := getOrSetId(w, r)
	data := id.GetAll()
	if data == nil {
		data := make(Session, 0)
		AllSessions[id] = data
	}
	return id
}

func init() {
	AllSessions = make(Sessions, 0)
}
