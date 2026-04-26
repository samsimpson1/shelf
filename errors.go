package main

import (
	"log"
	"net/http"
	"runtime"
)

// writeError writes a sanitised user-facing message to the response and logs
// the underlying error (with the calling site) for operators. Handlers must
// never embed err.Error() into the response body — wrapped error chains can
// leak filesystem paths, upstream URLs, or other internal state.
//
// Pass err = nil for input-validation errors that have no underlying cause.
func writeError(w http.ResponseWriter, status int, userMsg string, err error) {
	if err != nil {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			log.Printf("handler error at %s:%d: %s: %v", file, line, userMsg, err)
		} else {
			log.Printf("handler error: %s: %v", userMsg, err)
		}
	}
	http.Error(w, userMsg, status)
}
