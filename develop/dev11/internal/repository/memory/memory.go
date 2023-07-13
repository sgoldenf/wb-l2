package memory

import (
	"dev11/internal/repository"
	"dev11/pkg/model"
	"math/rand"
	"sync"
	"time"
)

// Repository is in-memory storage of Events where key is user_id and value is a map of events (key - event_id, value - Event).
// It's protected from concurrent read/write with sync.RWMutex.
// It uses *rand.Rand to generate event id.
type Repository struct {
	m          sync.RWMutex
	randomizer *rand.Rand
	data       map[uint64]map[uint64]*model.Event
}

// New creates an instance of repository and returns pointer to it
func New() *Repository {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Repository{randomizer: r, data: map[uint64]map[uint64]*model.Event{}}
}

// Create adds an Event to repository
func (r *Repository) Create(e *model.Event) (uint64, error) {
	e.ID = r.randomizer.Uint64()
	r.m.Lock()
	defer r.m.Unlock()
	if _, ok := r.data[e.UserID]; !ok {
		r.data[e.UserID] = make(map[uint64]*model.Event)
	}
	if _, ok := r.data[e.UserID][e.ID]; ok {
		return 0, repository.ErrDuplicateID
	}
	r.data[e.UserID][e.ID] = e
	return e.ID, nil
}

// Update changes an Event in repository
func (r *Repository) Update(e *model.Event) error {
	r.m.Lock()
	defer r.m.Unlock()
	if _, ok := r.data[e.UserID]; !ok {
		return repository.ErrUserNotFound
	}
	if _, ok := r.data[e.UserID][e.ID]; !ok {
		return repository.ErrEventNotFound
	}
	r.data[e.UserID][e.ID] = e
	return nil
}

// Delete removes an Event from repository
func (r *Repository) Delete(userID, id uint64) error {
	r.m.Lock()
	defer r.m.Unlock()
	if _, ok := r.data[userID]; !ok {
		return repository.ErrUserNotFound
	}
	if _, ok := r.data[userID][id]; !ok {
		return repository.ErrEventNotFound
	}
	delete(r.data[userID], id)
	return nil
}

// GetForDay returns a list of events for given day
func (r *Repository) GetForDay(userID uint64, t time.Time) ([]*model.Event, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	if _, ok := r.data[userID]; !ok {
		return nil, repository.ErrUserNotFound
	}
	events := []*model.Event{}
	for _, event := range r.data[userID] {
		if event.Date == t {
			events = append(events, event)
		}
	}
	return events, nil
}

// GetForWeek returns a list of events for a week starting from given day
func (r *Repository) GetForWeek(userID uint64, t time.Time) ([]*model.Event, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	if _, ok := r.data[userID]; !ok {
		return nil, repository.ErrUserNotFound
	}
	events := []*model.Event{}
	for _, event := range r.data[userID] {
		if event.Date == t || (event.Date.After(t) && event.Date.Before(t.AddDate(0, 0, 7))) {
			events = append(events, event)
		}
	}
	return events, nil
}

// GetForMonth returns a list of events for a month starting from given day
func (r *Repository) GetForMonth(userID uint64, t time.Time) ([]*model.Event, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	if _, ok := r.data[userID]; !ok {
		return nil, repository.ErrUserNotFound
	}
	events := []*model.Event{}
	for _, event := range r.data[userID] {
		if event.Date == t || (event.Date.After(t) && event.Date.Before(t.AddDate(0, 1, 0))) {
			events = append(events, event)
		}
	}
	return events, nil
}
