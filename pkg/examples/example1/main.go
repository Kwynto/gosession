package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/Kwynto/gosession"
)

func GetMd5(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	var html string = ""
	id := gosession.Start(&w, r)
	username := id.Get("username")

	if username != nil {
		html = `
		<html>
			<head>
				<title>Title</title>
			</head>
			<body>
				<p>
				You are authorized!<br>
				Your name is %s.<br>
				</p>
				<p>
				<a href="/logout">Logout?</a>
				</p>
				<p>
				<a href="/firstpage">First Page</a>
				</p>
				<p>
				<a href="/secondpage">Second Page</a>
				</p>
			</body>
		</html>
		`
		html = fmt.Sprintf(html, username)
	} else {
		html = `
		<html>
			<head>
				<title>Title</title>
			</head>
			<body>
				<p>
					<form action="/auth" method="post" class="form-horizontal">
						<input name="login" type="text" value="" placeholder="Login" required pattern="^[a-zA-Z0-9_-]+$">
						<input name="password" type="password" value="" placeholder="Password" required pattern="^[a-zA-Z0-9]+$">
						<button name="signin" type="submit">Auth button</button>
					</form>
				</p>
				<p>
					This is a test example of authorization on a web page.<br>
					Please enter any text as username and password
				</p>
			</body>
		</html>
		`
	}
	fmt.Fprint(w, html)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("login")
	password := r.FormValue("password")

	id := gosession.Start(&w, r)

	if username != "" && password != "" {
		pasHash := GetMd5(password)
		id.Set("username", username)
		id.Set("hash", pasHash)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func outHandler(w http.ResponseWriter, r *http.Request) {
	id := gosession.Start(&w, r)
	id.Destroy(&w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func firstHandler(w http.ResponseWriter, r *http.Request) {
	id := gosession.Start(&w, r)
	ses := id.GetAll()
	username := ses["username"]
	pasHash := ses["hash"]

	html := `
		<html>
			<head>
				<title>Title</title>
			</head>
			<body>
				<p>
				You are authorized!<br>
				</p>
				<p>
				Your name is %s.<br>
				Hash password is %s <br>
				</p>
				<p>
				<a href="/">Home Page</a>
				</p>
			</body>
		</html>
		`
	html = fmt.Sprintf(html, username, pasHash)
	fmt.Fprint(w, html)
}

func secondHandler(w http.ResponseWriter, r *http.Request) {
	id := gosession.Start(&w, r)

	html := `
		<html>
			<head>
				<title>Title</title>
			</head>
			<body>
				<p>
				You are authorized!<br>
				</p>
				<p>
				Your session ID is: %s.<br>
				</p>
				<p>
				<a href="/">Home Page</a>
				</p>
			</body>
		</html>
		`
	html = fmt.Sprintf(html, id)
	fmt.Fprint(w, html)
}

func main() {
	port := ":8080"

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/logout", outHandler)
	http.HandleFunc("/firstpage", firstHandler)
	http.HandleFunc("/secondpage", secondHandler)

	http.ListenAndServe(port, nil)
}
