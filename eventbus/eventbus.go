package eventbus

// EventBus offers the caller the interface to Subscribe or Publish to the EventBus, encapsulating its
// technical implementation.
type EventBus interface {
	EventSubscriber
	EventPublisher
}

type EventSubscriber interface {
	Subscribe(Stream, Subject) (chan Event, error)
}

type EventPublisher interface {
	Publish(Event) error
}

type (
	Subject string
	Stream  string
	Body    string
)

type Event struct {
	Stream  Stream
	Subject Subject
	Body    Body
}

type Subscription struct {
	EventBus chan Event
	Stream   Stream
	Subject  Subject
}
