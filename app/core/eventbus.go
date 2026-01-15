package core

// EventBus represents shared event bus
type EventBus struct {
	// This is a simplified implementation
	// In a real scenario, you would use a proper message bus
	subscribers map[string][]func(any)
}

// NewEventBus creates a new event bus instance
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]func(any)),
	}
}

// Subscribe subscribes to an event
func (eb *EventBus) Subscribe(event string, handler func(any)) {
	eb.subscribers[event] = append(eb.subscribers[event], handler)
}

// Publish publishes an event
func (eb *EventBus) Publish(event string, data any) {
	if handlers, exists := eb.subscribers[event]; exists {
		for _, handler := range handlers {
			handler(data)
		}
	}
}

// GetSubscribers returns the number of subscribers for an event
func (eb *EventBus) GetSubscribers(event string) int {
	return len(eb.subscribers[event])
}
