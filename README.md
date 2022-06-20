# GoSession
This is quick session for net/http in GoLang.
This package is perhaps the best implementation of the session mechanism, at least it tries to become one.

**Important note**
This package is designed to work with the standard net/http package and has not been tested with other http packages by the developer.

## Contents

- [GoSession](#gosession)
  - [Contents](#contents)
  - [What are sessions and why are they needed](#what-are-sessions-and-why-are-they-needed)
  - [How to connect GoSession](#how-to-connect-gosession)
  - [How to use GoSession](#how-to-use-gosession)
  - [Examples of using](#examples-of-using)
    - [Example 1](#example-1)
    - [Example 2](#example-2)
  - [About the package](#about-the-package)
  - [About the author](#about-the-author)

## What are sessions and why are they needed
A session on a site is a good method of identifying a site user.
A session is often used to authorize a user and retain their identity until the user closes the browser page.
While the user is working with the site, he saves cookies with a unique identifier, by this identifier one can distinguish one user from another and the server can store special data for a particular user.
User data received during the session period can be used for authorization, marketing and many other cases when it is necessary to collect, process and analyze data about a specific user.
A session is an efficient method of interacting with a user.

## How to connect GoSession
In your project folder, initialize the Go-module with the command
> go mod init your_app_name

Download and install GoSession
> go get github.com/Kwynto/gosession

Now you can add the GoSession package to your Go-code file, for example in `main.go`
```
import "github.com/Kwynto/gosession"
```

## How to use GoSession
To use the GoSession package, you need to import it into your code.
```
import "github.com/Kwynto/gosession"
```

All operations for working with sessions must be called from handlers.
Each time you start working with the session store, you need to call `gosession.Start(w *http.ResponseWriter, r *http.Request)`, since this function returns the identifier of the store and allows you to access the elements of the store through the identifier.
```
id := gosession.Start(&w, r)
```

You need to call the `gosession.Start(w *http.ResponseWriter, r *http.Request)` function from the handler
```
func rootHandler(w http.ResponseWriter, r *http.Request) {
  id := gosession.Start(&w, r) // Get the storage ID for a specific user

  html := "<html><head><title>Title</title></head><body>%s</body></html>"
  fmt.Fprintf(w, html, id)
}
```

Once you have a store ID, you can write variables to the store, read them, and delete them.

Recording is done using the `(id SessionId) Set(name string, value interface{})` method
```
id.Set("name variable", anyVariable)
```

In the handler it looks like this
```
func writeHandler(w http.ResponseWriter, r *http.Request) {
  name := "username"
  username := "JohnDow"

  id := gosession.Start(&w, r)
  id.Set(name, username)

  html := "<html><head><title>Title</title></head><body>OK</body></html>"
  fmt.Fprint(w, html)
}
```

Reading is done by `(id SessionId) Get(name string) interface{}` method for one variable
and the `(id SessionId) GetAll() Session` method to read all session variables
```
anyVariable := id.Get("name variable")
```

```
allVariables := id.GetAll()
```

In the handler it looks like this
```
func readHandler(w http.ResponseWriter, r *http.Request) {
  name := "username"
  var username interface{}

  id := gosession.Start(&w, r)
  username := id.Get(name) // Reading the "username" variable from the session for a specific user

  html := "<html><head><title>Title</title></head><body>%s</body></html>"
  fmt.Fprintf(w, html, username)
}
```

or so
```
func readHandler(w http.ResponseWriter, r *http.Request) {
  var tempStr string = ""

  id := gosession.Start(&w, r)
  allVariables := id.GetAll() // Reading the entire session for a specific client

  for i, v := range allVariables {
    tempStr = fmt.Sprint(tempStr, i, "=", v, "<br>")
  }
  html := "<html><head><title>Title</title></head><body>%s</body></html>"
  fmt.Fprintf(w, html, tempStr)
}
```

Removing an entry from a session of a specific client is carried out using the `(id SessionId) Remove(name string)` method
```
id.Remove("name variable")
```

In the handler it looks like this
```
func removeHandler(w http.ResponseWriter, r *http.Request) {
  id := gosession.Start(&w, r)
  id.Remove("name variable") // Removing a variable from a specific client session

  html := "<html><head><title>Title</title></head><body>OK</body></html>"
  fmt.Fprint(w, html)
}
```

Removing the entire session of a specific client is done using the `(id SessionId) Destroy(w *http.ResponseWriter)` method
```
id.Destroy(&w)
```

In the handler it looks like this
```
func destroyHandler(w http.ResponseWriter, r *http.Request) {
  id := gosession.Start(&w, r)
  id.Destroy(&w) // Deleting the entire session of a specific client

  html := "<html><head><title>Title</title></head><body>OK</body></html>"
  fmt.Fprint(w, html)
}
```

GoSession allows you to change its settings with the `SetSettings(setings GoSessionSetings)` function,
which is used outside of the handler, for example, inside the `main()` function
```
var mySetingsSession = gosession.GoSessionSetings{
  CookieName:    gosession.GOSESSION_COOKIE_NAME,
  Expiration:    gosession.GOSESSION_EXPIRATION,
  TimerCleaning: gosession.GOSESSION_TIMER_FOR_CLEANING,
}

gosession.SetSetings(mySetingsSession) // Setting session preferences
```

GoSession has 3 constants available for use
```
const (
  GOSESSION_COOKIE_NAME        string        = "SessionId" // Name for session cookies
  GOSESSION_EXPIRATION         int64         = 43_200      // Max age is 12 hours.
  GOSESSION_TIMER_FOR_CLEANING time.Duration = time.Hour   // The period of launch of the mechanism of cleaning from obsolete sessions
)
```

The remaining functions, types and variables in GoSession are auxiliary and are used only within the package.

## Examples of using

### Example 1:
*This is a simple authorization example and it shows the use of the write and read session variables functions, as well as deleting the entire session.*

Create an example folder and navigate to it.

Create a module for your application
> go mod init example1

Install GoSession
> go get github.com/Kwynto/gosession

Create a `main.go` file and save this code into it:
```
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
```

Run:
> go mod tidy

Start the server:
> go run .

Visit site
> http://localhost:8080/

### Example 2:


## About the package

GoSession has a description of its functionality in a `README.md` file and internal documentation.
GoSession is tested and has a performance check.
You can use the GoSession tests and documentation yourself.

Download the GoSession project to your computer:
> git clone https://github.com/Kwynto/gosession.git

Go to the project folder:
> cd ./gosession

**Check out the documentation**

Look at the documentation in two steps.
First, in the console, run:
> godoc -http=:8080

And then in your web browser navigate to the uri:
> http://localhost:8080

*The `godoc` utility may not be present in your Go build and you may need to install it
command `go get -v golang.org/x/tools/cmd/godoc`*

You can also use Go's standard functionality to view documentation in the console via `go doc`.
For example:
> go doc Start

If your IDE is good enough, then the documentation for functions and methods will be available from your code editor.

**Testing**

Run tests:
> go test -v

Run tests showing code coverage:
> go test -cover -v

You can view code coverage in detail in your web browser.
To do this, you need to sequentially execute two commands in the console:
> go test -coverprofile="coverage.out" -v
> go tool cover -html="coverage.out"

**Performance**

You can look at code performance tests:
> go test -benchmem -bench="." gosession.go gosession_test.go

*The slowest of all functions is `cleaningSessions()`, but this should not scare you, as it is a utility function and is rarely executed. This function does not affect the performance of the entire mechanism, it is only needed to clean up the storage from lost sessions.*

## About the author

The author of the project is Constantine Zavezeon (Kwynto).
You can contact the author by e-mail kwynto@mail.ru
The author accepts proposals for participation in open source projects,
as well as willing to accept job offers.