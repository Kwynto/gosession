package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Kwynto/gosession"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	var html string
	header := `
	<html>
		<head>
			<title>Title</title>
		</head>
		<body>
			<p>
			<a href="/">Home Page</a><br>
			<a href="/firstpage">First Page</a><br>
			<a href="/secondpage">Second Page</a><br>
			<a href="/thirdpage">Third Page</a><br>
			<a href="/fourthpage">Fourth Page</a><br>
			<a href="/fifthpage">Fifth Page</a><br>
			</p>
			<p>
			Website browsing history:<br>
	`

	footer := `
			</p>
		</body>
	</html>
	`

	id := gosession.Start(&w, r)
	transitions := id.Get("transitions")

	if transitions == nil {
		transitions = ""
	}
	transitions = fmt.Sprint(transitions, " ", r.RequestURI)
	id.Set("transitions", transitions)

	msg := fmt.Sprintf("%v", transitions)
	msg = strings.ReplaceAll(msg, " ", "<br>")
	html = fmt.Sprint(header, msg, footer)

	fmt.Fprint(w, html)
}

func favHandler(w http.ResponseWriter, r *http.Request) {
	// dummy
}

func main() {
	port := ":8080"

	http.HandleFunc("/favicon.ico", favHandler)
	http.HandleFunc("/", homeHandler)

	http.ListenAndServe(port, nil)
}
