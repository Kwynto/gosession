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
    - [Example 3](#example-3)
  - [About the package  (documentation, testing and benchmarking)](#about-the-package)
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
*This example shows a primitive way to collect information about user actions. You can collect any public user data, as well as track user actions, and then save and process this data.*

Create an example folder and navigate to it.

Create a module for your application
> go mod init example2

Install GoSession
> go get github.com/Kwynto/gosession

Create a `main.go` file and save this code into it:
```
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

```

Run:
> go mod tidy

Start the server:
> go run .

Visit site
> http://localhost:8080/

Now you can follow the links on this site and see how the site saves and shows your browsing history.

### Example 3:

*This example shows a simple, realistic site that uses the session mechanism.*

Create an example folder and navigate to it.

Create a module for your application
> go mod init example3

Install GoSession
> go get github.com/Kwynto/gosession

**Now you need to create 11 files.**

Create a `./cmd/web/main.go` file and save this code into it:
```
package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
)

// Creating an `application` structure to store the dependencies of the entire web application.
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	templateCache map[string]*template.Template
}

// Structure for configuration
type Config struct {
	Addr      string
	StaticDir string
}

func main() {
	// Reading flags from the application launch bar.
	cfg := new(Config)
	flag.StringVar(&cfg.Addr, "addr", ":8080", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	// Creating loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize the template cache.
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize the structure with the application dependencies.
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		templateCache: templateCache,
	}

	// Server structure with address, logger and routing
	srv := &http.Server{
		Addr:     cfg.Addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Starting the server
	app.infoLog.Printf("Start server: %s", cfg.Addr)
	err = srv.ListenAndServe()
	app.errorLog.Fatal(err)
}

```

Create a `./cmd/web/templates.go` file and save this code into it:
```
package main

import (
	"html/template"
	"path/filepath"
)

// Structure for the data template
type templateData struct {
	User        string
	Hash        string
	Cart        []string
	Transitions []string
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// We use the filepath.Glob function to get a slice of all file paths with the extension '.page.tmpl'.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// We iterate through the template file from each page.
	for _, page := range pages {
		// Extracting the final file name
		name := filepath.Base(page)

		// Processing the iterated template file.
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// We use the ParseGlob method to add all the wireframe templates.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// We use the ParseGlob method to add all auxiliary templates.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Adding the resulting set of templates to the cache using the page name
		cache[name] = ts
	}

	// We return the received map.
	return cache, nil
}

```

Create a `./cmd/web/routes.go` file and save this code into it:
```
package main

import (
	"net/http"
	"path/filepath"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

// Blocking direct access to the file system
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, _ := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

func (app *application) routes() *http.ServeMux {
	// Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/auth", app.authPage)
	mux.HandleFunc("/logout", app.outPage)
	mux.HandleFunc("/product1", app.buyPage)
	mux.HandleFunc("/product2", app.buyPage)
	mux.HandleFunc("/product3", app.buyPage)

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}

```

Create a `./cmd/web/handlers.go` file and save this code into it:
```
package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Kwynto/gosession"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	tD := &templateData{User: "", Hash: "", Cart: []string{""}, Transitions: []string{""}}

	id := gosession.Start(&w, r)

	transitions := id.Get("transitions")
	if transitions == nil {
		transitions = ""
	}
	transitions = fmt.Sprint(transitions, " ", r.RequestURI)
	id.Set("transitions", transitions)
	tStr := fmt.Sprintf("%v", transitions)
	tStrs := strings.Split(tStr, " ")
	tD.Transitions = tStrs

	cart := id.Get("cart")
	if cart == nil {
		cart = ""
		id.Set("cart", fmt.Sprint(cart))
		tD.Cart = []string{"There's nothing here yet."}
	} else {
		sCart := fmt.Sprint(cart)
		prods := strings.Split(sCart, " ")
		tD.Cart = prods
	}

	username := id.Get("username")
	if username != nil {
		tD.User = fmt.Sprint(username)
		app.render(w, r, "homeauth.page.tmpl", tD)
	} else {
		app.render(w, r, "home.page.tmpl", tD)
	}
}

func (app *application) authPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("login")
	password := r.FormValue("password")

	id := gosession.Start(&w, r)

	if username != "" && password != "" {
		pasHash := app.getMd5(password)
		id.Set("username", username)
		id.Set("hash", pasHash)
	}

	transitions := id.Get("transitions")
	if transitions == nil {
		transitions = ""
	}
	transitions = fmt.Sprint(transitions, " ", r.RequestURI)
	id.Set("transitions", transitions)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) outPage(w http.ResponseWriter, r *http.Request) {
	id := gosession.Start(&w, r)
	id.Remove("username")

	transitions := id.Get("transitions")
	if transitions == nil {
		transitions = ""
	}
	transitions = fmt.Sprint(transitions, " ", r.RequestURI)
	id.Set("transitions", transitions)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) buyPage(w http.ResponseWriter, r *http.Request) {
	tD := &templateData{User: "", Hash: "", Cart: []string{""}, Transitions: []string{""}}

	id := gosession.Start(&w, r)

	transitions := id.Get("transitions")
	if transitions == nil {
		transitions = ""
	}
	transitions = fmt.Sprint(transitions, " ", r.RequestURI)
	id.Set("transitions", transitions)
	tStr := fmt.Sprintf("%v", transitions)
	tStrs := strings.Split(tStr, " ")
	tD.Transitions = tStrs

	cart := id.Get("cart")
	if cart == nil {
		cart = ""
	}
	sCart := app.addProduct(fmt.Sprint(cart), app.convertProduct(r.RequestURI))
	id.Set("cart", sCart)
	prods := strings.Split(sCart, " ")
	tD.Cart = prods

	username := id.Get("username")
	if username != nil {
		tD.User = fmt.Sprint(username)
		app.render(w, r, "homeauth.page.tmpl", tD)
	} else {
		app.render(w, r, "home.page.tmpl", tD)
	}
}

```

Create a `./cmd/web/helpers.go` file and save this code into it:
```
package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
)

func (app *application) getMd5(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func (app *application) convertProduct(uri string) string {
	res := strings.Trim(uri, "/")
	res = strings.ToUpper(res)
	return res
}

func (app *application) addProduct(text string, plus string) string {
	var newText string = ""
	var newProd string = ""
	var isIt bool = false

	if text == "" {
		newText = fmt.Sprintf("%s=%d ", plus, 1)
		return newText
	}

	splitTest := strings.Split(text, " ")
	for _, val := range splitTest {
		if val != "" {
			splitVal := strings.Split(val, "=")
			prodName := splitVal[0]
			prodCount, _ := strconv.Atoi(splitVal[1])
			if prodName == plus {
				isIt = true
				prodCount += 1
				newProd = fmt.Sprintf("%s=%d", prodName, prodCount)
				newText = fmt.Sprintf("%s%s ", newText, newProd)
			} else {
				newText = fmt.Sprintf("%s%s ", newText, val)
			}
		}
	}

	if !isIt {
		newProd = fmt.Sprintf("%s=%d", plus, 1)
		newText = fmt.Sprintf("%s%s ", newText, newProd)
	}

	return newText
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// We extract the corresponding set of templates from the cache, depending on the page name
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("template %s not exist", name))
		return
	}

	// Render template files by passing dynamic data from the `td` variable.
	err := ts.Execute(w, td)
	if err != nil {
		app.serverError(w, err)
	}
}

```

Create a `./ui/html/base.layout.tmpl` file and save this code into it:
```
{{define "base"}}
<!doctype html>
<html lang='en'>
<head>
    <meta charset='utf-8'>
    <title>{{template "title" .}}</title>
    <link rel='stylesheet' href='/static/css/main.css'>
    <link rel='stylesheet' href='https://fonts.googleapis.com/css?family=Ubuntu+Mono:400,700'>
</head>
<body>
    <header>
        <h1><a href='/'>Home Page</a></h1>
    </header>
    <nav>
        <a href='/'>Home Page</a>
    </nav>
    <main>
        {{template "main" .}}
    </main>
    {{template "footer" .}}
    <script src="/static/js/main.js" type="text/javascript"></script>
</body>
</html>
{{end}}
```

Create a `./ui/html/footer.partial.tmpl` file and save this code into it:
```
{{define "footer"}}
<footer>This is the third example of a <strong>GoLang</strong> site using <strong>GoSession</strong>.</footer>
{{end}}
```

Create a `./ui/html/home.page.tmpl` file and save this code into it:
```
{{template "base" .}}

{{define "title"}}Home Page{{end}}

{{define "main"}}
	<table>
		<tr>
			<td rowspan="2" valign="top">
				<p>
					<h2>Authorization</h2>
					<form action="/auth" method="post" class="form-horizontal">
						<input name="login" type="text" value="" placeholder="Login" required pattern="^[a-zA-Z0-9_-]+$">
						<input name="password" type="password" value="" placeholder="Password" required pattern="^[a-zA-Z0-9]+$">
						<button name="signin" type="submit">Auth button</button>
					</form>
					<br><br>
				</p>
				<p>
					<h2>Links:</h2>
					<a href="/product1">Buy Product No. 1</a><br><br>
					<a href="/product2">Buy Product No. 2</a><br><br>
					<a href="/product3">Buy Product No. 3</a><br><br>
				</p>
			</td>
			<td width="250px">
				<p>
					<h2>Shopping cart</h2>
					{{ range $key, $pr := .Cart }}
						{{ $pr }}<br>
					{{end}}
				</p>
			</td>
		</tr>
		<tr>
			<td>
				<p>
					<h2>Browsing history:</h2>
					{{ range $key, $tr := .Transitions }}
						{{ $tr }}<br>
					{{end}}
				</p>
			</td>
		</tr>
	</table>
{{end}}
```

Create a `./ui/html/homeauth.page.tmpl` file and save this code into it:
```
{{template "base" .}}

{{define "title"}}Home Page{{end}}

{{define "main"}}
	<table>
		<tr>
			<td rowspan="2" valign="top">
				<p>
					<h3>You are logged in as: </h3>{{.User}} <a href="/logout">Log Out</a><br><br>
				</p>
				<p>
					<h2>Links:</h2>
					<a href="/product1">Buy Product No. 1</a><br><br>
					<a href="/product2">Buy Product No. 2</a><br><br>
					<a href="/product3">Buy Product No. 3</a><br><br>
				</p>
			</td>
			<td width="250px">
				<p>
					<h2>Shopping cart</h2>
					{{ range $key, $pr := .Cart }}
						{{ $pr }}<br>
					{{end}}
				</p>
			</td>
		</tr>
		<tr>
			<td>
				<p>
					<h2>Browsing history:</h2>
					{{ range $key, $tr := .Transitions }}
						{{ $tr }}<br>
					{{end}}
				</p>
			</td>
		</tr>
	</table>
{{end}}
```

Create a `./ui/static/css/main.css` file and save this code into it:
```
* {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
    font-size: 18px;
    font-family: "Ubuntu Mono", monospace;
}

html, body {
    height: 100%;
}

body {
    line-height: 1.5;
    background-color: #F1F3F6;
    color: #34495E;
    overflow-y: scroll;
}

header, nav, main, footer {
    padding: 2px calc((100% - 800px) / 2) 0;
}

main {
    margin-top: 54px;
    margin-bottom: 54px;
    min-height: calc(100vh - 345px);
    overflow: auto;
}

h1 a {
    font-size: 36px;
    font-weight: bold;
    background-image: url("/static/img/logo.png");
    background-repeat: no-repeat;
    background-position: 0px 0px;
    height: 36px;
    padding-left: 50px;
    position: relative;
}

h1 a:hover {
    text-decoration: none;
    color: #34495E;
}

h2 {
    font-size: 22px;
    margin-bottom: 36px;
    position: relative;
    top: -9px;
}

a {
    color: #62CB31;
    text-decoration: none;
}

a:hover {
    color: #4EB722;
    text-decoration: underline;
}

textarea, input:not([type="submit"]) {
    font-size: 18px;
    font-family: "Ubuntu Mono", monospace;
}

header {
    background-image: -webkit-linear-gradient(left, #34495e, #34495e 25%, #9b59b6 25%, #9b59b6 35%, #3498db 35%, #3498db 45%, #62cb31 45%, #62cb31 55%, #ffb606 55%, #ffb606 65%, #e67e22 65%, #e67e22 75%, #e74c3c 85%, #e74c3c 85%, #c0392b 85%, #c0392b 100%);
    background-image: -moz-linear-gradient(left, #34495e, #34495e 25%, #9b59b6 25%, #9b59b6 35%, #3498db 35%, #3498db 45%, #62cb31 45%, #62cb31 55%, #ffb606 55%, #ffb606 65%, #e67e22 65%, #e67e22 75%, #e74c3c 85%, #e74c3c 85%, #c0392b 85%, #c0392b 100%);
    background-image: -ms-linear-gradient(left, #34495e, #34495e 25%, #9b59b6 25%, #9b59b6 35%, #3498db 35%, #3498db 45%, #62cb31 45%, #62cb31 55%, #ffb606 55%, #ffb606 65%, #e67e22 65%, #e67e22 75%, #e74c3c 85%, #e74c3c 85%, #c0392b 85%, #c0392b 100%);
    background-image: linear-gradient(to right, #34495e, #34495e 25%, #9b59b6 25%, #9b59b6 35%, #3498db 35%, #3498db 45%, #62cb31 45%, #62cb31 55%, #ffb606 55%, #ffb606 65%, #e67e22 65%, #e67e22 75%, #e74c3c 85%, #e74c3c 85%, #c0392b 85%, #c0392b 100%);
    background-size: 100% 6px;
    background-repeat: no-repeat;
    border-bottom: 1px solid #E4E5E7;
    overflow: auto;
    padding-top: 33px;
    padding-bottom: 27px;
    text-align: center;
}

header a {
    color: #34495E;
    text-decoration: none;
}

nav {
    border-bottom: 1px solid #E4E5E7;
    padding-top: 17px;
    padding-bottom: 15px;
    background: #F7F9FA;
    height: 60px;
    color: #6A6C6F;
}

nav a {
    margin-right: 1.5em;
    display: inline-block;
}

nav form {
    display: inline-block;
    margin-left: 1.5em;
}

nav div {
    width: 50%;
    float: left;
}

nav div:last-child {
    text-align: right;
}

nav div:last-child a {
    margin-left: 1.5em;
    margin-right: 0;
}

nav a.live {
    color: #34495E;
    cursor: default;
}

nav a.live:hover {
    text-decoration: none;
}

nav a.live:after {
    content: '';
    display: block;
    position: relative;
    left: calc(50% - 7px);
    top: 9px;
    width: 14px;
    height: 14px;
    background: #F7F9FA;
    border-left: 1px solid #E4E5E7;
    border-bottom: 1px solid #E4E5E7;
    -moz-transform: rotate(45deg);
    -webkit-transform: rotate(-45deg);
}

form div {
    margin-bottom: 18px;
}

form div:last-child {
    border-top: 1px dashed #E4E5E7;
}

form input[type="radio"] {
    position: relative;
    top: 2px;
    margin-left: 18px;
}

form input[type="text"], form input[type="password"], form input[type="email"] {
    padding: 0.75em 18px;
    width: 100%;
}

form input[type=text], form input[type="password"], form input[type="email"], textarea {
    color: #6A6C6F;
    background: #FFFFFF;
    border: 1px solid #E4E5E7;
    border-radius: 3px;
}

form label {
    display: inline-block;
    margin-bottom: 9px;
}

.error {
    color: #C0392B;
    font-weight: bold;
    display: block;
}

.error + textarea, .error + input {
    border-color: #C0392B !important;
    border-width: 2px !important;
}

textarea {
    padding: 18px;
    width: 100%;
    height: 266px;
}

button {
    background-color: #4CAF50;
    border: none;
    color: white;
    padding: 15px 32px;
    text-align: center;
    text-decoration: none;
    display: inline-block;
    font-size: 16px;
    margin: 4px 2px;
    cursor: pointer;
    width: 100%;
}

button:hover {
    background-color: #5865f4;
    color: white;
}

.snippet {
    background-color: #FFFFFF;
    border: 1px solid #E4E5E7;
    border-radius: 3px;
}

.snippet pre {
    padding: 18px;
    border-top: 1px solid #E4E5E7;
    border-bottom: 1px solid #E4E5E7;
}

.snippet .metadata {
    background-color: #F7F9FA;
    color: #6A6C6F;
    padding: 0.75em 18px;
    overflow: auto;
}

.snippet .metadata span {
    float: right;
}

.snippet .metadata strong {
    color: #34495E;
}

.snippet .metadata time {
    display: inline-block;
}

.snippet .metadata time:first-child {
    float: left;
}

.snippet .metadata time:last-child {
    float: right;
}

div.flash {
    color: #FFFFFF;
    font-weight: bold;
    background-color: #34495E;
    padding: 18px;
    margin-bottom: 36px;
    text-align: center;
}

div.error {
    color: #FFFFFF;
    background-color: #C0392B;
    padding: 18px;
    margin-bottom: 36px;
    font-weight: bold;
    text-align: center;
}

table {
    background: white;
    border: 1px solid #E4E5E7;
    border-collapse: collapse;
    width: 100%;
}

td, th {
    text-align: left;
    padding: 9px 18px;
}

th:last-child, td:last-child {
    text-align: right;
    color: #6A6C6F;
}

tr {
    border-bottom: 1px solid #E4E5E7;
}

tr:nth-child(2n) {
    background-color: #F7F9FA;
}

footer {
    border-top: 1px solid #E4E5E7;
    padding-top: 17px;
    padding-bottom: 15px;
    background: #F7F9FA;
    height: 60px;
    color: #6A6C6F;
    text-align: center;
}

```

Create a `./ui/static/js/main.js` file and save this code into it:
```
var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}
```

Run:
> go mod tidy

Start the server:
> go run ./cmd/web/

Visit site
> http://localhost:8080/

Now you can follow the links on this site.

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
