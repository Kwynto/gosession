package gosession

import "net/http"

type Sessions map[string]interface{}

var AllSessions Sessions

func generateId() {

}

func Put(id string, v interface{}) {

}

func Get(id string) (string, interface{}) {

}

func SessionStart(r *http.Request) (string, error) {

	nil
}

func init() {
	AllSessions = make(Sessions, 10)

}
