package eventbus

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EventBusTestSuite struct {
	suite.Suite
	EventBus RegioBus
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

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		incomingEvent := <-incomingEvents
		s.Equal(orderEvent, incomingEvent)
		s.NotEqual(invoiceEvent.Stream, incomingEvent)
		wg.Done()
	}()

	err = s.EventBus.Publish(orderEvent)
	s.NoError(err)
	err = s.EventBus.Publish(invoiceEvent)
	s.NoError(err)

	wg.Wait()
}

func TestEventBusTestSuite(t *testing.T) {
	suite.Run(t, new(EventBusTestSuite))
}
