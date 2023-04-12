# IncomingEvents package

This code provides a Go package for handling incoming events. It includes the following functions:

## Types

- `Event`: represents an incoming event and includes the event name, variable, and data.
- `RoutingCriteria`: represents the criteria used for routing events to the appropriate handlers.
- `IncomingEvents`: represents the set of handlers for incoming events and includes a map of routing criteria to handlers, a lock to handle concurrent access, a map of aggregated events, and an aggregation size.

## Functions

- `NewIncomingEvents(aggregateSize int)`: creates a new IncomingEvents object with the specified aggregation size.
- `AddHandler(criteria RoutingCriteria, handler EventHandler) error`: adds a handler for the specified routing criteria.
- `triggerHandlers(criteria RoutingCriteria, event Event) error`: triggers the handlers for the specified routing criteria with the specified event.
- `RemoveHandler(criteria RoutingCriteria, handler EventHandler) error`: removes the specified handler for the specified routing criteria.
- `ParseEvent(eventString string, filterCriteria *RoutingCriteria) error`: parses the incoming event string and triggers the appropriate handlers based on the specified routing criteria and aggregation size.

This package allows you to handle incoming events by creating and managing handlers that are triggered based on the routing criteria defined in the `RoutingCriteria` type. You can also aggregate events based on the criteria defined in `RoutingCriteria` and the aggregation size provided in `NewIncomingEvents`.
