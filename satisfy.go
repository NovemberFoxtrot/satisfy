package main

import (
	"net/http"
	"net/http/httptest"
)

type ModifierMiddleware struct {
	handler http.Handler
}

func (m *ModifierMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rec := httptest.NewRecorder()

	m.handler.ServeHTTP(rec, r)

	for k, v := range rec.Header() {
		w.Header()[k] = v
	}

	w.Header().Set("X-Hey-Hey-Hey", "Yup")
	w.WriteHeader(418)
	w.Write([]byte("Middleware says hello again..."))
	w.Write(rec.Body.Bytes())
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Success!!"))
}

func main() {
	mid := &ModifierMiddleware{http.HandlerFunc(myHandler)}

	println("on 8080")

	http.ListenAndServe(":8080", mid)
}
