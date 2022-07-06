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

	id := gosession.StartSecure(&w, r)

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

	id := gosession.StartSecure(&w, r)

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
	id := gosession.StartSecure(&w, r)
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

	id := gosession.StartSecure(&w, r)

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
