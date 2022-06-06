package gosession

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"
)

const (
	GOSESSION_COOKIE_NAME        string        = "SessionId"
	GOSESSION_EXPIRATION         int64         = 43_200    // Max age is 12 hours.
	GOSESSION_TIMER_FOR_CLEANING time.Duration = time.Hour // 1 hour
)

type SessionId string

type Session map[string]interface{}

type internalSession struct {
	expiration int64
	data       Session
}

type serverSessions map[SessionId]internalSession

var allSessions serverSessions

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

func cleaningSessions() {
	presently := time.Now().Unix()
	for id, ses := range allSessions {
		if ses.expiration < presently {
			delete(allSessions, id)
		}
	}
	time.AfterFunc(GOSESSION_TIMER_FOR_CLEANING, cleaningSessions)
}

func (id SessionId) Set(name string, value interface{}) {
	ses := allSessions[id]
	ses.data[name] = value
	allSessions[id] = ses
}

func (id SessionId) GetAll() Session {
	return allSessions[id].data
}

func (id SessionId) GetOne(name string) interface{} {
	ses := allSessions[id]
	return ses.data[name]
}

func (id SessionId) RemoveSession(w *http.ResponseWriter) {
	delete(allSessions, id)
	deleteCookie(w)
}

func (id SessionId) RemoveValue(name string) {
	ses := allSessions[id]
	delete(ses.data, name)
	allSessions[id] = ses
}

func Start(w *http.ResponseWriter, r *http.Request) SessionId {
	id := getOrSetCookie(w, r)
	ses := allSessions[id]
	if ses.data == nil {
		ses.data = make(Session, 0)
	}
	presently := time.Now().Unix()
	ses.expiration = presently + GOSESSION_EXPIRATION
	allSessions[id] = ses
	return id
}

// Package initialization

func init() {
	allSessions = make(serverSessions, 0)
	time.AfterFunc(GOSESSION_TIMER_FOR_CLEANING, cleaningSessions)
	fmt.Println("GoSessions initialized")
}
