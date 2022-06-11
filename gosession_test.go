package gosession

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"
)

const (
	GOSESSION_TESTING_ITER int = 10
)

var result SessionId

// ----------------------------
// Helper functions for testing
// ----------------------------

func rootHandler(w http.ResponseWriter, r *http.Request) {
	Start(&w, r)
	html := `
	<html>
		<head>
			<title>Title</title>
		</head>
		<body>
			<form action="/auth" method="post" class="form-horizontal">
				<input name="login" type="text" value="" placeholder="Login" required pattern="^[a-zA-Z0-9_-]+$">
				<input name="password" type="password" value="" placeholder="Password" required pattern="^[a-zA-Z0-9]+$">
				<button name="signin" type="submit">Auth button</button>
			</form>
		</body>
	</html>
	`
	fmt.Fprint(w, html)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	sesid := Start(&w, r)
	name := sesid.GetOne("username")
	password := sesid.GetOne("password")
	floatnumber := sesid.GetOne("float")
	intnumber := sesid.GetOne("number")
	construct := sesid.GetOne("construct")
	allses := sesid.GetAll()

	cleaningSessions()

	html := "<html><head><title>Title</title></head><body>This is a test!<br>Username: %s<br>Password: %s<br>%v<br>%v<br>%v<br>%v<br></body></html>"
	fmt.Fprintf(w, html, name, password, floatnumber, intnumber, construct, allses)

	sesid.RemoveValue("construct")
	sesid.RemoveValue("dest")

	sesid.RemoveSession(&w)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	sesid := Start(&w, r)
	name := r.PostFormValue("login")
	password := r.PostFormValue("password")
	sesid.Set("username", name)
	sesid.Set("password", password)
	sesid.Set("float", 3.14)
	sesid.Set("number", 13)
	tstruct := struct {
		name string
		pas  string
		fnum float64
		inum int
	}{
		name: name,
		pas:  password,
		fnum: 2.2,
		inum: 15,
	}
	sesid.Set("construct", tstruct)
	html := "<html><head><title>Title</title></head><body>%s<br><a href='/test'>Test session</a></body></html>"
	fmt.Fprintf(w, html, name)
}

func realServer() {
	SetSetings(GoSessionSetings{CookieName: "goSessionID", Expiration: 40, TimerCleaning: time.Second * 90})
	PORT := ":8001"
	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/", rootHandler)
	http.ListenAndServe(PORT, nil)
}

// --------------
// Test functions
// --------------

func Test_generateId(t *testing.T) {
	testVar := make(map[int]SessionId)
	for i := 0; i < GOSESSION_TESTING_ITER; i++ {
		testVar[i] = generateId()
	}
	for _, v1 := range testVar {
		count := 0
		for _, v2 := range testVar {
			if bytes.Equal([]byte(v1), []byte(v2)) {
				count++
			}
		}
		if count > 1 {
			t.Error("Error generating unique identifier.")
		}
	}
}

func Test_realServer(t *testing.T) {
	go realServer()
	for i := 0; i < 100; i++ {
		time.Sleep(time.Second)
	}
	// for i := 0; i < GOSESSION_TESTING_ITER; i++ {
	// 	req, err := http.NewRequest("GET", "/get", nil)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	rr := httptest.NewRecorder()
	// 	handler := http.HandlerFunc(getData)
	// 	handler.ServeHTTP(rr, req)
	// }
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
