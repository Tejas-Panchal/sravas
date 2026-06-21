package middleware

var validEventTypes = map[string]bool{
	"view":           true,
	"watch_time":     true,
	"like":           true,
	"comment":        true,
	"traffic_source": true,
}

// ValidateEventType checks whether the given event type is supported
func ValidateEventType(eventType string) string {
	if !validEventTypes[eventType] {
		return "unsupported event type (valid: view, watch_time, like, comment, traffic_source)"
	}
	return ""
}
