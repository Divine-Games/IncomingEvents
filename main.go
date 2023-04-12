how do you fix this error in the code

panic: runtime error: invalid memory address or nil pointer dereference
[signal 0xc0000005 code=0x0 addr=0x28 pc=0xf809c2]

goroutine 1 [running]:
github.com/Divine-Games/IncomingEvents.(*IncomingEvents).ParseEvent(0xc00006e3c0, {0xfa3e05, 0xef5ab9}, 0x0)
        E:/Gaea/Golang/pkg/mod/github.com/!divine-!games/!incoming!events@v0.0.3/main.go:127 +0x3e2

package IncomingEvents

import (
	"errors"
	"reflect"
	"strings"
	"sync"
)

type EventHandler func(event Event)

type Event struct {
    Name     string
    Variable string
    Data     string
}

type RoutingCriteria struct {
    EventType string
    EventData string
}

type IncomingEvents struct {
    handlers      map[RoutingCriteria][]EventHandler
    lock          sync.Mutex
    aggregated    map[RoutingCriteria]*Event
    aggregateSize int
}

func NewIncomingEvents(aggregateSize int) *IncomingEvents {
    return &IncomingEvents{
        handlers:      make(map[RoutingCriteria][]EventHandler),
        lock:          sync.Mutex{},
        aggregated:    make(map[RoutingCriteria]*Event, 0),
        aggregateSize: aggregateSize,
    }
}

func (ie *IncomingEvents) AddHandler(criteria RoutingCriteria, handler EventHandler) error {
    ie.lock.Lock()
    defer ie.lock.Unlock()

    if criteria.EventType == "" {
        return errors.New("EventType cannot be empty")
    }

    handlers := ie.handlers[criteria]
    handlers = append(handlers, handler)
    ie.handlers[criteria] = handlers

    return nil
}

func (ie *IncomingEvents) triggerHandlers(criteria RoutingCriteria, event Event) error {
    ie.lock.Lock()
    defer ie.lock.Unlock()

    handlers, ok := ie.handlers[criteria]
    if !ok {
        return errors.New("no handlers registered for criteria")
    }

    for _, handler := range handlers {
        go handler(event)
    }

    return nil
}

func (ie *IncomingEvents) RemoveHandler(criteria RoutingCriteria, handler EventHandler) error {
    ie.lock.Lock()
    defer ie.lock.Unlock()

    if criteria.EventType == "" {
        return errors.New("EventType cannot be empty")
    }

    handlers, ok := ie.handlers[criteria]
    if !ok {
        return errors.New("no handlers registered for criteria")
    }

    for i, h := range handlers {
        if reflect.ValueOf(h).Pointer() == reflect.ValueOf(handler).Pointer() {
            handlers = append(handlers[:i], handlers[i+1:]...)
            ie.handlers[criteria] = handlers
            return nil
        }
    }

    return errors.New("handler not found")
}

func (ie *IncomingEvents) ParseEvent(eventString string, filterCriteria *RoutingCriteria) error {
    eventParts := strings.SplitN(eventString, ":", 3)
    if len(eventParts) != 3 {
        return errors.New("invalid event format")
    }

    eventType := eventParts[0]
    if eventType == "" {
        return errors.New("EventType cannot be empty")
    }

    eventVariable := eventParts[1]
    eventData := strings.TrimSpace(eventParts[2])

    event := Event{Name: eventType, Variable: eventVariable, Data: eventData}

    // Define the routing criteria based on the event type and data
    criteria := RoutingCriteria{EventType: eventType, EventData: eventData}

    ie.lock.Lock()
    defer ie.lock.Unlock()

	

	if ie.aggregated[criteria] == nil {
		ie.aggregated[criteria] = &Event{}
	}

    // Apply filter criteria if provided
    if filterCriteria != nil {
        if filterCriteria.EventType != "" && filterCriteria.EventType != eventType {
            return nil
        }
        if filterCriteria.EventData != "" && filterCriteria.EventData != eventData {
            return nil
        }
    }

    // Aggregate the event data
    if len(ie.aggregated[criteria].Data) == 0 {
        ie.aggregated[criteria] = &event
    } else {
        ie.aggregated[criteria].Data += ", " + event.Data
    }

    // Trigger the handlers with the new event if it matches the filter criteria
    handlers, ok := ie.handlers[criteria]
    if ok {
        for _, handler := range handlers {
            if filterCriteria == nil ||
                (filterCriteria.EventType == "" || filterCriteria.EventType == eventType) &&
                (filterCriteria.EventData == "" || filterCriteria.EventData == eventData) {
                go handler(event)
            }
        }
    }

    // If the aggregation size is reached, trigger the handlers with the aggregated event if it matches the filter criteria
    if len(ie.aggregated[criteria].Data) >= ie.aggregateSize {
        aggregatedEvent := *ie.aggregated[criteria]
        ie.aggregated[criteria] = &Event{}

        if ok {
            for _, handler := range handlers {
                if filterCriteria == nil ||
                    (filterCriteria.EventType == "" || filterCriteria.EventType == eventType) &&
                    (filterCriteria.EventData == "" || filterCriteria.EventData == eventData) {
                    go handler(aggregatedEvent)
                }
            }
        }
    }

    return nil
}
