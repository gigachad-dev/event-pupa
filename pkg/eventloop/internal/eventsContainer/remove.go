package eventsContainer

import (
	"gitlab.com/YSX/eventloop/pkg/eventloop/event"
)

func (el *eventsList) removeEvent(e event.Interface) {
	el.removeByTypes(e)
	el.removeByPriority(e)
	el.removeByTrigger(e)
}

func (el *eventsList) removeTrigger(trig string) {
	triggerInfo := el.eventsByCriteria[TRIGGER][trig]
	for _, ev := range triggerInfo.data {
		el.removeEvent(ev)
	}
	// Удаляем триггер после всех удалений
	if len(triggerInfo.data) == 0 && len(triggerInfo.middlewareEvents) == 0 {
		delete(el.eventsByCriteria[TRIGGER], trig)
	}
}

func (el *eventsList) removeByTypes(e event.Interface) {
	for _, t := range e.GetTypes() {
		el._removeHelper(TYPE, string(t), e.GetUUID())
	}
}

func (el *eventsList) removeByTrigger(e event.Interface) {
	triggerName := e.GetTriggerName()
	el._removeHelper(TRIGGER, triggerName, e.GetUUID())
}

func (el *eventsList) removeByPriority(e event.Interface) {
	priority := e.GetPriorityString()

	el._removeHelper(PRIORITY, priority, e.GetUUID())
}

func (el *eventsList) _removeHelper(criteriaName criteriaName, criteria string, uuid string) {
	if priorityInfo, ok := el.eventsByCriteria[criteriaName][criteria]; ok {
		delete(priorityInfo.data, uuid)

		if len(priorityInfo.data) == 0 {
			delete(el.eventsByCriteria[criteriaName], criteria)
		}
	}
}
