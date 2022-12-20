package event

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"

	"eventloop/pkg/eventloop/event/after"
	"eventloop/pkg/eventloop/event/interval"
	"eventloop/pkg/eventloop/event/once"
	"eventloop/pkg/eventloop/event/subscriber"
	"github.com/google/uuid"
)

// event - обычное событие, которое может иметь свойства других событий (одноразовых, интервальных, зависимых)

type EventArgs struct {
	TriggerName string
	Priority    int
	IsOnce      bool
	Fun         EventFunc

	IntervalTime time.Duration
	after.DateAfter
}

type event struct {
	id          uuid.UUID
	triggerName string
	priority    int
	fun         EventFunc
	result      string

	mx sync.Mutex

	subscriber subscriber.Interface
	interval   interval.Interface
	once       once.Interface
	after      after.Interface
}

type EventFunc func(ctx context.Context) string

func NewEvent(args EventArgs) (Interface, error) {
	if args.Fun == nil {
		return nil, errors.New("no function, please add")
	}

	newEvent := &event{id: uuid.New(),
		fun:         args.Fun,
		triggerName: args.TriggerName,
		priority:    args.Priority}

	if args.IsOnce {
		newEvent.once = once.NewOnce()
	}
	if args.IntervalTime.String() != "0s" {
		newEvent.interval = interval.NewIntervalEvent(args.IntervalTime)
	}
	if args.DateAfter != (after.DateAfter{}) {
		newEvent.after = after.New(args.DateAfter)
	}

	return newEvent, nil
}

func (ev *event) GetID() uuid.UUID {
	return ev.id
}

func (ev *event) GetTriggerName() string {
	return ev.triggerName
}

func (ev *event) GetPriority() int {
	return ev.priority
}

func (ev *event) RunFunction(ctx context.Context) {
	ev.result = ev.fun(ctx)

	//Отправка сообщений, подписанным на это событие, событиям
	listener := ev.Subscriber()
	if listenerChannels := listener.GetChannels(); len(listenerChannels) > 0 {
		listener.LockMutex()
		for _, chnl := range listenerChannels {
			chnl <- 1
		}
		listener.UnlockMutex()
	}
}

// Subscriber
func (ev *event) Subscriber() subscriber.Interface {
	if ev.subscriber == nil {
		ev.subscriber = subscriber.NewSubscriber()
	}
	return ev.subscriber
}

func (ev *event) Interval() (interval.Interface, error) {
	return getSubInterface(ev.interval, "it is not an interval event")
}

func (ev *event) Once() (once.Interface, error) {
	return getSubInterface(ev.once, "it is not an once event")
}

func (ev *event) After() (after.Interface, error) {
	return getSubInterface(ev.after, "it is not an after event")
}

func getSubInterface[T any](i T, errMsg string) (T, error) {
	var nilRes T
	if reflect.ValueOf(&i).Elem().IsZero() {
		return nilRes, errors.New(errMsg)
	}
	return i, nil
}
