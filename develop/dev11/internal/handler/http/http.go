package http

import (
	"dev11/internal/controller/event"
	"errors"
	"net/http"
)

var (
	errInvalidEventID = errors.New("invalid event id")
	errInvalidUserID  = errors.New("invalid user id")
	errEmptyTitle     = errors.New("empty title")
	errInvalidDate    = errors.New("invalid date")
)

// Handler processes HTTP requests
type Handler struct {
	ctrl *event.Controller
}

// New creates Handler instance with provided Controller and returns pointer to it
func New(ctrl *event.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// PostCreateEvent handles POST HTTP Request to add Event to calendar
func (h *Handler) PostCreateEvent(w http.ResponseWriter, req *http.Request) {
	e, err := parseEvent(req)
	if err != nil && !errors.Is(err, errInvalidEventID) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.ctrl.Create(e)
	if err != nil {
		if errors.Is(err, event.ErrDuplicateID) {
			writeError(w, http.StatusServiceUnavailable, http.StatusText(http.StatusServiceUnavailable))
		} else {
			writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
		return
	}
	writeResponseJSON(w, http.StatusCreated, map[string]interface{}{"result": id})
}

// PostUpdateEvent handles POST HTTP Request to change Event in calendar
func (h *Handler) PostUpdateEvent(w http.ResponseWriter, req *http.Request) {
	e, err := parseEvent(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = h.ctrl.Update(e)
	if err != nil {
		if errors.Is(err, event.ErrUserNotFound) || errors.Is(err, event.ErrEventNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
		return
	}
	writeResponseJSON(w, http.StatusOK, map[string]interface{}{"result": "successfully updated"})
}

// PostDeleteEvent handles POST HTTP Request to remove Event from calendar
func (h *Handler) PostDeleteEvent(w http.ResponseWriter, req *http.Request) {
	userID, err := parseUserID(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	eventID, err := parseEventID(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.ctrl.Delete(userID, eventID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	writeResponseJSON(w, http.StatusOK, map[string]interface{}{"result": "successfully deleted"})
}

// GetEventsForDay handles GET HTTP Request for an events occuring at given day
func (h *Handler) GetEventsForDay(w http.ResponseWriter, req *http.Request) {
	userID, err := parseUserID(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	date, err := parseDate(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	events, err := h.ctrl.GetForDay(userID, date)
	if err != nil {
		if errors.Is(err, event.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
		return
	}

	writeResponseJSON(w, http.StatusOK, map[string]interface{}{"result": events})
}

// GetEventsForWeek handles GET HTTP Request for an events occuring in a week starting from given day
func (h *Handler) GetEventsForWeek(w http.ResponseWriter, req *http.Request) {
	userID, err := parseUserID(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	date, err := parseDate(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	events, err := h.ctrl.GetForWeek(userID, date)
	if err != nil {
		if errors.Is(err, event.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
		return
	}

	writeResponseJSON(w, http.StatusOK, map[string]interface{}{"result": events})
}

// GetEventsForMonth handles GET HTTP Request for an events occuring in a month starting from given day
func (h *Handler) GetEventsForMonth(w http.ResponseWriter, req *http.Request) {
	userID, err := parseUserID(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	date, err := parseDate(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	events, err := h.ctrl.GetForMonth(userID, date)
	if err != nil {
		if errors.Is(err, event.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
		return
	}

	writeResponseJSON(w, http.StatusOK, map[string]interface{}{"result": events})
}
