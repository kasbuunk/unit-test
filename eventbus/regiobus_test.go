package eventbus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type EventBusTestSuite struct {
	suite.Suite
	EventBus *RegioBus
}

func (s *EventBusTestSuite) SetupTest() {
	regioBus := New([]string{"ORDER", "INVOICE", "USER"})
	s.EventBus = regioBus
}

func (s *EventBusTestSuite) TestPubSub() {
	// Init some participants in the event bus. Some subscribers and publishers.
	orderStream := Stream("ORDER")
	invoiceStream := Stream("INVOICE")
	incomingEvents, err := s.EventBus.Subscribe(orderStream, "*")
	s.NoError(err)

	orderEvent := Event{
		Stream:  orderStream,
		Subject: "order placed",
	}
	invoiceEvent := Event{
		Stream:  invoiceStream,
		Subject: "invoice paid",
	}
	go func() {
		incomingEvent := <-incomingEvents
		s.Equal(orderEvent, incomingEvent)
		s.NotEqual(invoiceEvent.Stream, incomingEvent)
	}()

	err = s.EventBus.Publish(orderEvent)
	s.NoError(err)
	time.Sleep(20 * time.Millisecond)
}

func TestEventBusTestSuite(t *testing.T) {
	suite.Run(t, new(EventBusTestSuite))
}
