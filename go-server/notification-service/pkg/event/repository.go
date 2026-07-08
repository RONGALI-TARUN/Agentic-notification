package event

import "sync"

type Repository interface {
	Save(event Event) error
	List() ([]Event, error)
}

type InMemoryRepository struct {
	mu     sync.Mutex
	events []Event
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{}
}

func (r *InMemoryRepository) Save(evt Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, evt)
	return nil
}

func (r *InMemoryRepository) List() ([]Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Event, len(r.events))
	copy(out, r.events)
	return out, nil
}
