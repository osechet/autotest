package autotest

import (
	"fmt"
	"reflect"
)

var Any = "autotest.Any"

type EventID uint32

type event struct {
	id   EventID
	args []string
}

var events = make(chan event)

// QueueClosedError is returned when the queue has been unexpectedly closed.
type QueueClosedError struct{}

func (QueueClosedError) Error() string {
	return "The queue has been closed"
}

// UnexpectedEventError is returned when the received event is not the expected one.
type UnexpectedEventError struct {
	expectedID   EventID
	expectedArgs []string
	receivedID   EventID
	receivedArgs []string
}

func (e UnexpectedEventError) Error() string {
	return fmt.Sprintf("Unexpected event: received {%d, %v} - expected {%d, %v}", e.receivedID, e.receivedArgs, e.expectedID, e.expectedArgs)
}

// Clear clears the event queue.
func Clear() {
L:
	for {
		select {
		case <-events:
		default:
			break L
		}
	}
}

// Expect is called to ensure the next event is the given one.
func Expect(id EventID, args ...string) (EventID, []string, error) {
	event, ok := <-events
	if !ok {
		return 0, nil, QueueClosedError{}
	}
	if args == nil {
		args = make([]string, 0)
	}
	if event.id != id {
		return 0, nil, UnexpectedEventError{id, args, event.id, event.args}
	}
	if len(args) == 0 || args[0] != Any {
		if !reflect.DeepEqual(event.args, args) {
			return 0, nil, UnexpectedEventError{id, args, event.id, event.args}
		}
	}
	return event.id, event.args, nil
}

// Skip skips the next event.
func Skip() (EventID, []string, error) {
	event, ok := <-events
	if !ok {
		return 0, nil, QueueClosedError{}
	}
	return event.id, event.args, nil
}
