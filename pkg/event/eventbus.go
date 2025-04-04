package event

const (
	EventLinkVisited = "link.visited"
)

type Event struct {
	Type string
	Data any
}

type EventBus struct {
	Bus chan Event
}

func NewEventBus() *EventBus {
	return &EventBus{
		Bus: make(chan Event),
	}
}

func (eb *EventBus) Publish(event Event) {
	eb.Bus <- event
}

func (eb *EventBus) Subsribe() <-chan Event {
	return eb.Bus
}
