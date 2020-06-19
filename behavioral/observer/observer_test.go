package observer

import "testing"

func TestObserver(t *testing.T) {

	publisher := NewPublisher()

	publisher.Attach(&ConcreteObserver{})
	publisher.Attach(&ConcreteObserver{})
	publisher.Attach(&ConcreteObserver{})

	publisher.SetState("New State...")

	publisher.Notify()
}
