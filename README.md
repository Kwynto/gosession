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
  - [About the package](#about-the-package)
  - [About the author](#about-the-author)

## What are sessions and why are they needed

## How to connect GoSession
In your project folder, initialize the Go-module with the command
> go mod init your_app_name

Download and install GoSession
> go get github.com/Kwynto/gosession

Now you can add the GoSession package to your Go-code file, for example in `main.go`
> import "github.com/Kwynto/gosession"

## How to use GoSession
To use the GoSession package, you need to import it into your code.
> import "github.com/Kwynto/gosession"

All operations for working with sessions must be called from handlers.
Each time you start working with the session store, you need to call `gosession.Start(w *http.ResponseWriter, r *http.Request)`, since this function returns the identifier of the store and allows you to access the elements of the store through the identifier.
>  id := gosession.Start(&w, r)

You need to call the `gosession.gosessionStart(w *http.ResponseWriter, r *http.Request)` function from the handler
>func rootHandler(w http.ResponseWriter, r *http.Request) {
>  id := gosession.Start(&w, r) // Get the storage ID for a specific user
> 
>  html := "\<html\>\<head\>\<title\>Title\<\/title\>\<\/head\>\<body\>%s\<\/body\>\<\/html\>"
>  fmt.Fprintf(w, html, id)
>}

Once you have a store ID, you can write variables to the store, read them, and delete them.

Recording is done using the `(id SessionId) Set(name string, value interface{})` method
> id.Set("name variable", anyVariable)

In the handler it looks like this
>func writeHandler(w http.ResponseWriter, r *http.Request) {
>  name := "username"
>  username := "JohnDow"
> 
>  id := gosession.Start(&w, r)
>  id.Set(name, username)
> 
>  html := "\<html\>\<head\>\<title\>Title\<\/title\>\<\/head\>\<body\>OK\<\/body\>\<\/html\>"
>  fmt.Fprint(w, html)
>}

Reading is done by `(id SessionId) Get(name string) interface{}` method for one variable
and the `(id SessionId) GetAll() Session` method to read all session variables
> anyVariable := id.Get("name variable")

> allVariables := id.GetAll()

In the handler it looks like this
>func readHandler(w http.ResponseWriter, r *http.Request) {
>  name := "username"
>  var username interface{}
> 
>  id := gosession.Start(&w, r)
>  username := id.Get(name) // Reading the "username" variable from the session for a specific user
> 
>  html := "\<html\>\<head\>\<title\>Title\<\/title\>\<\/head\>\<body\>%s\<\/body\>\<\/html\>"
>  fmt.Fprintf(w, html, username)
>}

or so
>func readHandler(w http.ResponseWriter, r *http.Request) {
>  var tempStr string = ""
> 
>  id := gosession.Start(&w, r)
>  allVariables := id.GetAll() // Reading the entire session for a specific client
> 
>  for i, v := range allVariables {
>    tempStr = fmt.Sprint(tempStr, i, "=", v, "\<br\>")
>  }
>  html := "\<html\>\<head\>\<title\>Title\<\/title\>\<\/head\>\<body\>%s\<\/body\>\<\/html\>"
>  fmt.Fprintf(w, html, tempStr)
>}

Удалление записи из сессии конкретного клиента осуществляется методом `(id SessionId) RemoveValue(name string)`
> id.RemoveValue("name variable")

In the handler it looks like this
>func removeHandler(w http.ResponseWriter, r *http.Request) {
>  id := gosession.Start(&w, r)
>  id.RemoveValue("name variable") // Удаление переменной из сессии конкретного клиента
>  
>  html := "\<html\>\<head\>\<title\>Title\<\/title\>\<\/head\>\<body\>OK\<\/body\>\<\/html\>"
>  fmt.Fprint(w, html)
>}

Удаление всей сессии конкретного клиента осуществляется методом `(id SessionId) RemoveSession(w *http.ResponseWriter)`
> id.RemoveSession(&w)

In the handler it looks like this
>func removeHandler(w http.ResponseWriter, r *http.Request) {
>  id := gosession.Start(&w, r)
>  id.RemoveSession(&w) // Удаление всей сессии конкретного клиента
>  
>  html := "\<html\>\<head\>\<title\>Title\<\/title\>\<\/head\>\<body\>OK\<\/body\>\<\/html\>"
>  fmt.Fprint(w, html)
>}

GoSession позволяет изменять свои настройки функцией `gosession.SetSetings(setings GoSessionSetings)`, которая используется за пределами хендлера, например, внутри функции `main()`
>var mySetingsSession = gosession.GoSessionSetings{
>  CookieName:    gosession.GOSESSION_COOKIE_NAME,
>  Expiration:    gosession.GOSESSION_EXPIRATION,
>  TimerCleaning: gosession.GOSESSION_TIMER_FOR_CLEANING,
>}
>
>gosession.SetSetings(mySetingsSession) // Установка настроек сессии

GoSession имеет три константы доступные для использования
>const (
>	GOSESSION_COOKIE_NAME        string        = "SessionId" // Name for session cookies
>	GOSESSION_EXPIRATION         int64         = 43_200      // Max age is 12 hours.
>	GOSESSION_TIMER_FOR_CLEANING time.Duration = time.Hour   // The period of launch of the mechanism of cleaning from obsolete sessions
>)

Остальные функции, типы и переменные в GoSession являются вспомогательными и используются только внутри пакета.

## Examples of using

## About the package

## About the author

