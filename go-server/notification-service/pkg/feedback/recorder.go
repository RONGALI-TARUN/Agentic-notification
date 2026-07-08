package feedback

import (
	"sync"
	"time"
)

type Recorder interface {
	Record(feedback EventFeedback) error
	List() ([]EventFeedback, error)
}

type InMemoryRecorder struct {
	mu       sync.Mutex
	feedback []EventFeedback
}

func NewInMemoryRecorder() *InMemoryRecorder {
	return &InMemoryRecorder{}
}

func (r *InMemoryRecorder) Record(feedback EventFeedback) error {
	if feedback.CreatedAt == "" {
		feedback.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.feedback = append(r.feedback, feedback)
	return nil
}

func (r *InMemoryRecorder) List() ([]EventFeedback, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]EventFeedback, len(r.feedback))
	copy(out, r.feedback)
	return out, nil
}
