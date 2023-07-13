package event

import (
	"dev11/internal/repository"
	"dev11/pkg/model"
	"errors"
	"time"
)

// Errors from Repository
var (
	ErrUserNotFound  = errors.New("user not found")
	ErrEventNotFound = errors.New("event not found")
	ErrDuplicateID   = errors.New("duplicate event id")
)

type eventRepository interface {
	Create(e *model.Event) (uint64, error)
	Update(e *model.Event) error
	Delete(userID, id uint64) error
	GetForDay(userID uint64, t time.Time) ([]*model.Event, error)
	GetForWeek(userID uint64, t time.Time) ([]*model.Event, error)
	GetForMonth(userID uint64, t time.Time) ([]*model.Event, error)
}

// Controller contains an instance of repository and provides its methods to client
type Controller struct {
	repo eventRepository
}

// New creates an instance of Controller provided with repository and returns pointer to it
func New(repo eventRepository) *Controller {
	return &Controller{repo: repo}
}

// Create adds an Event to repository
func (c *Controller) Create(e *model.Event) (uint64, error) {
	return c.repo.Create(e)
}

// Update changes an Event from repository
func (c *Controller) Update(e *model.Event) error {
	err := c.repo.Update(e)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFound
		}
		if errors.Is(err, repository.ErrEventNotFound) {
			return ErrEventNotFound
		}
	}
	return err
}

// Delete removes an Event from repository
func (c *Controller) Delete(userID, id uint64) error {
	err := c.repo.Delete(userID, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFound
		}
		if errors.Is(err, repository.ErrEventNotFound) {
			return ErrEventNotFound
		}
	}
	return err
}

// GetForDay returns a list of events for given day
func (c *Controller) GetForDay(userID uint64, t time.Time) ([]*model.Event, error) {
	events, err := c.repo.GetForDay(userID, t)
	if err != nil && errors.Is(err, repository.ErrUserNotFound) {
		return nil, ErrUserNotFound
	}
	return events, err
}

// GetForWeek returns a list of events for a week starting from given day
func (c *Controller) GetForWeek(userID uint64, t time.Time) ([]*model.Event, error) {
	events, err := c.repo.GetForWeek(userID, t)
	if err != nil && errors.Is(err, repository.ErrUserNotFound) {
		return nil, ErrUserNotFound
	}
	return events, err
}

// GetForMonth returns a list of events for a month starting from given day
func (c *Controller) GetForMonth(userID uint64, t time.Time) ([]*model.Event, error) {
	events, err := c.repo.GetForMonth(userID, t)
	if err != nil && errors.Is(err, repository.ErrUserNotFound) {
		return nil, ErrUserNotFound
	}
	return events, err
}
