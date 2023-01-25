package eventslist

import (
	"eventloop/pkg/eventloop/event"
)

type Interface interface {
	TriggerName(triggerName string) Priority
	RemoveEventByUUIDs(uuids ...string) []string
}

type Priority interface {
	Priority(priority int) *EventsByUUIDString
	Len() int
	GetSortedPriorityNums() (keys []int)
	IsDisabled() bool
	SetIsDisabled(b bool)
	GetAllEvents() (result []string, err error)
}

type EventID interface {
	List() EventsByUUIDString
	EventID(eventID string) event.Interface
	AddEvent(newEvent event.Interface)
}
