// This is quick session for net/http in golang.
package gosession

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"sync"
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

var block sync.RWMutex

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

// The deleteCookie(w) function deletes the session cookie
func deleteCookie(w *http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   setingsSession.CookieName,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(*w, cookie)
}

// The cleaningSessions() function periodically cleans up the server's session storage
func cleaningSessions() {
	presently := time.Now().Unix()
	block.Lock()
	for id, ses := range allSessions {
		if ses.expiration < presently {
			delete(allSessions, id)
		}
	}
	block.Unlock()
	// log.Println("Session storage has been serviced.")
	time.AfterFunc(setingsSession.TimerCleaning, cleaningSessions)
}

// The writeS() method safely writes data to the session store
func (id SessionId) writeS(iSes internalSession) {
	block.Lock()
	allSessions[id] = iSes
	block.Unlock()
}

// The readS() method safely reads data from the session store.
func (id SessionId) readS() (internalSession, bool) {
	block.RLock()
	defer block.RUnlock()
	ses, ok := allSessions[id]
	if !ok {
		return internalSession{}, false
	}
	return ses, true
}

// The destroyS() method safely deletes the entire session from the store.
func (id SessionId) destroyS() {
	block.Lock()
	delete(allSessions, id)
	block.Unlock()
}

// The deleteS() method safely deletes one client variable from the session by its name
// name - session variable name
func (id SessionId) deleteS(name string) {
	block.Lock()
	ses, ok := allSessions[id]
	if ok {
		delete(ses.data, name)
		allSessions[id] = ses
	}
	block.Unlock()
}

// The Set(name, value) SessionId-method to set the client variable to be stored in the session system.
// name - session variable name.
// value - directly variable in session.
func (id SessionId) Set(name string, value interface{}) {
	ses, ok := id.readS()
	if ok {
		ses.data[name] = value
		id.writeS(ses)
	}
}

// The GetAll() SessionId-method to get all client variables from the session system
func (id SessionId) GetAll() Session {
	ses, _ := id.readS()
	return ses.data
}

// The Get(name) SessionId-method to get a specific client variable from the session system.
// name - session variable name
func (id SessionId) Get(name string) interface{} {
	ses, _ := id.readS()
	return ses.data[name]
}

// The Destroy(w) SessionId-method to remove the entire client session
func (id SessionId) Destroy(w *http.ResponseWriter) {
	id.destroyS()
	deleteCookie(w)
}

// The Remove(name) SessionId-method to remove one client variable from the session by its name
func (id SessionId) Remove(name string) {
	id.deleteS(name)
}

// The SetSetings(settings) sets new settings for the session mechanism.
// setings - gosession.GoSessionSetings public type variable for setting new session settings
func SetSetings(setings GoSessionSetings) {
	setingsSession = setings
}

// The Start(w, r) function starts the session and returns the SessionId to the handler for further use of the session mechanism.
// This function must be run at the very beginning of the http.Handler
func Start(w *http.ResponseWriter, r *http.Request) SessionId {
	id := getOrSetCookie(w, r)
	ses, ok := id.readS()
	if !ok {
		ses.data = make(Session, 0)
	}
	presently := time.Now().Unix()
	ses.expiration = presently + setingsSession.Expiration
	id.writeS(ses)
	return id
}

// The StartSecure(w, r) function starts the session or changes the session ID and sets new cookie to the client.
// This function must be run at the very beginning of the http.Handler
func StartSecure(w *http.ResponseWriter, r *http.Request) SessionId {
	id := getOrSetCookie(w, r)
	ses, ok := id.readS()
	if !ok {
		ses.data = make(Session, 0)
		presently := time.Now().Unix()
		ses.expiration = presently + setingsSession.Expiration
		id.writeS(ses)
		return id
	} else {
		id.destroyS()
		id = generateId()
		cookie := &http.Cookie{
			Name:   setingsSession.CookieName,
			Value:  string(id),
			MaxAge: 0,
		}
		http.SetCookie(*w, cookie)
		presently := time.Now().Unix()
		ses.expiration = presently + setingsSession.Expiration
		id.writeS(ses)
		return id
	}
}

// Package initialization
func init() {
	time.AfterFunc(setingsSession.TimerCleaning, cleaningSessions)
	log.Println("GoSessions initialized")
}
