package peerforge

type EventType string

type Event struct {
	// ID is a Version 4 UUID
	ID string `json:"id"`

	// Version starting from 0 (zero).
	Version *int `json:"version,string"`

	// Source is a unique source string identifier
	Source string `json:"source"`

	// Type is a general event type
	Type EventType `json:"type"`
}

func NewEvent(t EventType, id string, version int, source string) *Event {
	return &Event{
		ID:      id,
		Version: &version,
		Source:  source,
		Type:    t,
	}
}
