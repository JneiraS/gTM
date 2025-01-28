package api

import (
	// "text/template"
	"net/http"
)

// StartServer will start a server on port 8080 and map the "/" path to the
// task handler's ServeTasks method.
func StartServer() {

	http.ListenAndServe(":8080", nil)
}
