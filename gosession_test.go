package gosession

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	GOSESSION_TESTING_ITER int = 10000
)

var result SessionId

// --------------
// Test functions
// --------------

func Test_generateId(t *testing.T) {
	testVar := make(map[int]SessionId)
	for i := 0; i < GOSESSION_TESTING_ITER; i++ {
		testVar[i] = generateId() // calling the tested function
	}
	for _, v1 := range testVar {
		count := 0
		for _, v2 := range testVar {
			if bytes.Equal([]byte(v1), []byte(v2)) {
				count++
			}
		}
		// work check
		if count > 1 {
			t.Error("Error generating unique identifier.")
		}
	}
}

func Test_getOrSetCookie(t *testing.T) {
	for i := 0; i < GOSESSION_TESTING_ITER; i++ {
		var ctrlId SessionId
		handler := func(w http.ResponseWriter, r *http.Request) {
			sesid := getOrSetCookie(&w, r) // calling the tested function
			ctrlId = sesid
			io.WriteString(w, string(sesid))
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		handler(w, r)

		status := w.Code
		// work check
		if status != http.StatusOK {
			t.Errorf("Handler returned %v", status)
		}

		cookies := w.Result().Cookies()
		noErr := false
		for _, v := range cookies {
			if v.Name == setingsSession.CookieName && v.Value == string(ctrlId) {
				noErr = true
			}
		}
		// work check
		if !noErr {
			t.Error("the server returned an invalid ID")
		}
	}

	for i := 0; i < GOSESSION_TESTING_ITER; i++ {
		var ctrlId SessionId
		handler := func(w http.ResponseWriter, r *http.Request) {
			sesid := getOrSetCookie(&w, r) // calling the tested function
			ctrlId = sesid
			io.WriteString(w, string(sesid))
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		clientId := generateId()
		cookie := &http.Cookie{
			Name:   setingsSession.CookieName,
			Value:  string(clientId),
			MaxAge: 0,
		}
		r.AddCookie(cookie)
		handler(w, r)

		status := w.Code
		// work check
		if status != http.StatusOK {
			t.Errorf("Handler returned %v", status)
		}
		// work check
		if ctrlId != clientId {
			t.Error("server received invalid id")
		}
	}
}

func Test_deleteCookie(t *testing.T) {
	for i := 0; i < GOSESSION_TESTING_ITER; i++ {
		handler := func(w http.ResponseWriter, r *http.Request) {
			deleteCookie(&w) // calling the tested function
			io.WriteString(w, "<html><head><title>Title</title></head><body>Body</body></html>")
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		clientId := generateId()
		cookie := &http.Cookie{
			Name:   setingsSession.CookieName,
			Value:  string(clientId),
			MaxAge: 0,
		}
		r.AddCookie(cookie)
		handler(w, r)

		status := w.Code
		// work check
		if status != http.StatusOK {
			t.Errorf("Handler returned %v", status)
		}

		cookies := w.Result().Cookies()
		noErr := true
		for _, v := range cookies {
			if v.Name == setingsSession.CookieName && v.Value == string(clientId) {
				noErr = false
			}
		}
		// work check
		if !noErr {
			t.Error("the server did not delete the session cookie")
		}
	}
}

func Test_cleaningSessions(t *testing.T) {
	rand.Seed(time.Now().Unix())
	var falseInd int
	var trueInd int
	for i := 0; i < GOSESSION_TESTING_ITER; i++ {
		falseInd = rand.Intn(75)
		trueInd = rand.Intn(50) + falseInd

		for id := range allSessions {
			delete(allSessions, id)
		}

		for fi := 0; fi < falseInd; fi++ {
			allSessions[generateId()] = internalSession{
				expiration: 0,
				data:       make(Session),
			}
		}

		for ti := 0; ti < trueInd; ti++ {
			allSessions[generateId()] = internalSession{
				expiration: time.Now().Unix() + setingsSession.Expiration,
				data:       make(Session),
			}
		}

		cleaningSessions() // calling the tested function
		// work check
		if len(allSessions) != trueInd {
			t.Error("The number of correct entries does not match.")
		}
	}
}

func Test_Set(t *testing.T) {
	var value interface{}
	rand.Seed(time.Now().Unix())

	for i := 0; i < GOSESSION_TESTING_ITER; i++ {
		id := generateId()
		allSessions[id] = internalSession{
			expiration: time.Now().Unix() + setingsSession.Expiration,
			data:       make(Session),
		}

		name := "test variable"
		switch rand.Intn(3) {
		case 0:
			value = true
		case 1:
			value = fmt.Sprintf("test string %d", rand.Intn(100))
		case 2:
			value = rand.Intn(100)
		case 3:
			value = rand.Float64()
		}

		id.Set(name, value) // calling the tested function
		// work check
		if allSessions[id].data[name] != value {
			t.Error("Failed to write variable to session storage.")
		}
	}
}

func Test_GetAll(t *testing.T) {
	var value interface{}
	var name string
	rand.Seed(time.Now().Unix())

	for i := 0; i < GOSESSION_TESTING_ITER; i++ {
		id := generateId()
		data := make(Session)
		count := rand.Intn(20) + 1
		for ic := 0; ic < count; ic++ {
			name = fmt.Sprintf("test name  %d", rand.Intn(100))
			switch rand.Intn(3) {
			case 0:
				value = true
			case 1:
				value = fmt.Sprintf("test string %d", rand.Intn(100))
			case 2:
				value = rand.Intn(100)
			case 3:
				value = rand.Float64()
			}
			data[name] = value
		}
		allSessions[id] = internalSession{
			expiration: time.Now().Unix() + setingsSession.Expiration,
			data:       data,
		}

		ses := id.GetAll() // calling the tested function
		// work check
		for iname, v := range ses {
			if v != data[iname] {
				t.Error("Incorrect data received from session variable storage")
			}
		}
	}
}

func Test_Get(t *testing.T) {
	var value interface{}
	var name string
	rand.Seed(time.Now().Unix())

	for i := 0; i < GOSESSION_TESTING_ITER; i++ {
		id := generateId()
		data := make(Session)

		name = "test name"
		switch rand.Intn(3) {
		case 0:
			value = true
		case 1:
			value = fmt.Sprintf("test string %d", rand.Intn(100))
		case 2:
			value = rand.Intn(100)
		case 3:
			value = rand.Float64()
		}
		data[name] = value
		allSessions[id] = internalSession{
			expiration: time.Now().Unix() + setingsSession.Expiration,
			data:       data,
		}

		getedValue := id.Get(name) // calling the tested function
		// work check
		if getedValue != value {
			t.Error("Incorrect data received from session variable storage")
		}
	}
}

func Test_RemoveSession(t *testing.T) {
	var hid SessionId

	for i := 0; i < GOSESSION_TESTING_ITER; i++ {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		id := generateId()
		cookie := &http.Cookie{
			Name:   setingsSession.CookieName,
			Value:  string(id),
			MaxAge: 0,
		}
		r.AddCookie(cookie)

		data := make(Session)
		name := "test name"
		value := "test value"
		data[name] = value
		allSessions[id] = internalSession{
			expiration: time.Now().Unix() + setingsSession.Expiration,
			data:       data,
		}

		handler := func(w http.ResponseWriter, r *http.Request) {
			hid = Start(&w, r)
			hid.RemoveSession(&w) // calling the tested function
			io.WriteString(w, "<html><head><title>Title</title></head><body>Body</body></html>")
		}
		handler(w, r)

		status := w.Code
		// work check
		if status != http.StatusOK {
			t.Errorf("Handler returned status: %v", status)
		}

		// work check
		if id != hid {
			t.Error("ID mismatch")
		}

		cookies := w.Result().Cookies()
		noErr := true
		for _, v := range cookies {
			if v.Name == setingsSession.CookieName && v.Value == string(id) {
				noErr = false
			}
		}
		// work check
		if !noErr {
			t.Error("The server did not delete the session cookie")
		}

		// work check
		if allSessions[id].data != nil {
			t.Error("Session has not been deleted.")
		}
	}
}

func Test_RemoveValue(t *testing.T) {
	var value interface{}
	var name string
	rand.Seed(time.Now().Unix())

	for i := 0; i < GOSESSION_TESTING_ITER; i++ {
		id := generateId()
		data := make(Session)

		name = "test name"
		switch rand.Intn(3) {
		case 0:
			value = true
		case 1:
			value = fmt.Sprintf("test string %d", rand.Intn(100))
		case 2:
			value = rand.Intn(100)
		case 3:
			value = rand.Float64()
		}
		data[name] = value
		allSessions[id] = internalSession{
			expiration: time.Now().Unix() + setingsSession.Expiration,
			data:       data,
		}

		id.RemoveValue(name) // calling the tested function
		// work check
		if allSessions[id].data[name] == value {
			t.Error("Failed to change settings")
		}
	}
}

func Test_SetSetings(t *testing.T) {
	var test_setingsSession1 = GoSessionSetings{
		CookieName:    "test_name",
		Expiration:    int64(rand.Intn(86_400)),
		TimerCleaning: time.Minute,
	}
	var test_setingsSession2 = GoSessionSetings{
		CookieName:    GOSESSION_COOKIE_NAME,
		Expiration:    GOSESSION_EXPIRATION,
		TimerCleaning: GOSESSION_TIMER_FOR_CLEANING,
	}

	for i := 0; i < GOSESSION_TESTING_ITER; i++ {
		SetSetings(test_setingsSession1) // calling the tested function
		// work check
		if test_setingsSession1 != setingsSession {
			t.Error("Failed to change settings.")
		}
		SetSetings(test_setingsSession2) // calling the tested function
		// work check
		if test_setingsSession2 != setingsSession {
			t.Error("Failed to change settings.")
		}
	}
}

func Test_Start(t *testing.T) {
	for i := 0; i < GOSESSION_TESTING_ITER; i++ {
		var id1 SessionId
		var id2 SessionId
		handler1 := func(w http.ResponseWriter, r *http.Request) {
			id1 = Start(&w, r) // calling the tested function
			io.WriteString(w, "<html><head><title>Title</title></head><body>Body</body></html>")
		}
		handler2 := func(w http.ResponseWriter, r *http.Request) {
			id2 = Start(&w, r) // calling the tested function
			io.WriteString(w, "<html><head><title>Title</title></head><body>Body</body></html>")
		}

		w1 := httptest.NewRecorder()
		r1 := httptest.NewRequest("GET", "/", nil)
		handler1(w1, r1)

		status := w1.Code
		// work check
		if status != http.StatusOK {
			t.Errorf("Handler returned status: %v", status)
		}

		cookies := w1.Result().Cookies()
		var cookie *http.Cookie
		for _, v := range cookies {
			cookie = v
		}
		w2 := httptest.NewRecorder()
		r2 := httptest.NewRequest("GET", "/", nil)
		r2.AddCookie(cookie)
		handler2(w2, r2)

		status = w2.Code
		// work check
		if status != http.StatusOK {
			t.Errorf("Handler returned status: %v", status)
		}

		// work check
		if id1 != id2 {
			t.Errorf("Server and client IDs are not equal:\n server: %v\n client: %v\n", id1, id2)
		}
	}
}

// ---------------------------------
// Helper functions for benchmarking
// ---------------------------------

// ----------------------
// Functions benchmarking
// ----------------------

func Benchmark_generateId(b *testing.B) {
	var r SessionId
	for i := 0; i < b.N; i++ {
		r = generateId()
	}
	result = r
}
