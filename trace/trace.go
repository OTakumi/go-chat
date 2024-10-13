package trace

// Trace is the interface that describes the behavior of a trace.
type Tracer interface {
	Trace(...interface{})
}
