package main

import (
  "net/http"
  "log"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
  log.Printf("internal server error!: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
  
  writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}
func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request)  {
  log.Printf("forbidden!: %s path: %s error: %s", r.Method, r.URL.Path)
  
  writeJSONError(w, http.StatusForbidden, "forbidden")

}
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
  log.Printf("bad request error!: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
  
  writeJSONError(w, http.StatusBadRequest, err.Error())
}
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
  log.Printf("not found error!: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
  
  writeJSONError(w, http.StatusNotFound, "resource not found")
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("unauthorized error: %s path: %s error:", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

