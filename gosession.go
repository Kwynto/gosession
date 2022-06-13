// This is quick session for net/http in golang.
package gosession

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"
)

const (
	GOSESSION_COOKIE_NAME        string        = "SessionId" // Name for session cookies
	GOSESSION_EXPIRATION         int64         = 43_200      // Max age is 12 hours.
	GOSESSION_TIMER_FOR_CLEANING time.Duration = time.Hour   // The period of launch of the mechanism of cleaning from obsolete sessions
)

// The SessionId type is the session identifier
type SessionId string

// The Session type contains variables defined for session storage for each client.
type Session map[string]interface{}

// The internalSession type is the internal server representation of the session
type internalSession struct {
	expiration int64
	data       Session
}

// The serverSessions type is intended to describe all sessions of all client connections
type serverSessions map[SessionId]internalSession

// The GoSessionSetings type describes the settings for the session system
type GoSessionSetings struct {
	CookieName    string
	Expiration    int64
	TimerCleaning time.Duration
}

// The allSessions variable stores all sessions of all clients
var allSessions serverSessions = make(serverSessions, 0)

// Session mechanism settings variable
var setingsSession = GoSessionSetings{
	CookieName:    GOSESSION_COOKIE_NAME,
	Expiration:    GOSESSION_EXPIRATION,
	TimerCleaning: GOSESSION_TIMER_FOR_CLEANING,
}

// The generateId() generates a new session id in a random, cryptographically secure manner
func generateId() SessionId {
	b := make([]byte, 32)
	rand.Read(b)
	return SessionId(fmt.Sprintf("%x", b))
}

// The getOrSetCookie(w, r) gets the session id from the cookie, or creates a new one if it can't get
func getOrSetCookie(w *http.ResponseWriter, r *http.Request) SessionId {
	data, err := r.Cookie(setingsSession.CookieName)
	if err != nil {
		id := generateId()
		cookie := &http.Cookie{
			Name:   setingsSession.CookieName,
			Value:  string(id),
			MaxAge: 0,
		}
		http.SetCookie(*w, cookie)
		return id
	}
	return SessionId(data.Value)
}

// The deleteCookie(w) deletes the session cookie
func deleteCookie(w *http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   setingsSession.CookieName,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(*w, cookie)
}

// The cleaningSessions() periodically cleans up the server's session storage
func cleaningSessions() {
	presently := time.Now().Unix()
	for id, ses := range allSessions {
		if ses.expiration < presently {
			delete(allSessions, id)
		}
	}
	time.AfterFunc(setingsSession.TimerCleaning, cleaningSessions)
}

// The Set(name, value) SessionId-method to set the client variable to be stored in the session system
// name - session variable name
// value - directly variable in session
func (id SessionId) Set(name string, value interface{}) {
	ses := allSessions[id]
	ses.data[name] = value
	allSessions[id] = ses
}

// The GetAll() SessionId-method to get all client variables from the session system
func (id SessionId) GetAll() Session {
	return allSessions[id].data
}

// The Get(name) SessionId-method to get a specific client variable from the session system
// name - session variable name
func (id SessionId) Get(name string) interface{} {
	ses := allSessions[id]
	return ses.data[name]
}

// The RemoveSession(w) SessionId-method to remove the entire client session
func (id SessionId) RemoveSession(w *http.ResponseWriter) {
	delete(allSessions, id)
	deleteCookie(w)
}

// The RemoveValue(name) SessionId-method to remove one client variable from the session by its name
func (id SessionId) RemoveValue(name string) {
	ses := allSessions[id]
	delete(ses.data, name)
	allSessions[id] = ses
}

// The SetSetings(settings) sets new settings for the session mechanism
// setings - gosession.GoSessionSetings public type variable for setting new session settings
func SetSetings(setings GoSessionSetings) {
	setingsSession = setings
}

// The Start(w, r) starts the session and returns the SessionId to the handler for further use of the session mechanism.
// This function must be run at the very beginning of the http.Handler
func Start(w *http.ResponseWriter, r *http.Request) SessionId {
	id := getOrSetCookie(w, r)
	ses := allSessions[id]
	if ses.data == nil {
		ses.data = make(Session, 0)
	}
	presently := time.Now().Unix()
	ses.expiration = presently + setingsSession.Expiration
	allSessions[id] = ses
	return id
}

// Package initialization
func init() {
	time.AfterFunc(setingsSession.TimerCleaning, cleaningSessions)
	fmt.Println("GoSessions initialized")
}
