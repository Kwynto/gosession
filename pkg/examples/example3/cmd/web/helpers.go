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
