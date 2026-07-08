package feedback

import "testing"

func TestInMemoryRecorderStoresFeedback(t *testing.T) {
	recorder := NewInMemoryRecorder()
	feedback := EventFeedback{ID: "fb-1", EventID: "evt-1", Action: "opened"}

	if err := recorder.Record(feedback); err != nil {
		t.Fatalf("record feedback: %v", err)
	}

	stored, err := recorder.List()
	if err != nil {
		t.Fatalf("list feedback: %v", err)
	}
	if len(stored) != 1 {
		t.Fatalf("expected 1 feedback item, got %d", len(stored))
	}
	if stored[0].Action != "opened" {
		t.Fatalf("expected action opened, got %s", stored[0].Action)
	}
}
