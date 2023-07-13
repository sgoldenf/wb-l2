package http

import (
	"dev11/pkg/model"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

func parseEvent(req *http.Request) (e *model.Event, err error) {
	userID, err := parseUserID(req)
	if err != nil {
		return
	}

	title := req.FormValue("title")
	if title == "" {
		return nil, errEmptyTitle
	}

	date, err := parseDate(req)
	if err != nil {
		return
	}

	id, err := parseEventID(req)
	e = &model.Event{ID: id, UserID: userID, Title: title, Date: date}
	return
}

func parseEventID(req *http.Request) (uint64, error) {
	idValue := req.FormValue("id")
	id, err := strconv.ParseUint(idValue, 10, 64)
	if err != nil {
		return id, errInvalidEventID
	}
	return id, nil
}

func parseUserID(req *http.Request) (uint64, error) {
	idValue := req.FormValue("user_id")
	id, err := strconv.ParseUint(idValue, 10, 64)
	if err != nil {
		return id, errInvalidUserID
	}
	return id, nil
}

func parseDate(req *http.Request) (time.Time, error) {
	dateValue := req.FormValue("date")
	date, err := time.Parse("2006-01-02", dateValue)
	if err != nil {
		return date, errInvalidDate
	}
	return date, nil
}

func writeResponseJSON(w http.ResponseWriter, code int, data interface{}) {
	resp, _ := json.Marshal(data)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	resp := map[string]string{"error": msg}
	writeResponseJSON(w, code, resp)
}

// Post is a middleware for POST HTTP methods
func (h *Handler) Post(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			writeError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Get is a middleware for GET HTTP methods
func (h *Handler) Get(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			writeError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Log is a middleware for logging a request
func (h *Handler) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}
