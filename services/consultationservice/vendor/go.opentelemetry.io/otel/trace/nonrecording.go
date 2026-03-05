


package trace 




type nonRecordingSpan struct {
	noopSpan

	sc SpanContext
}


func (s nonRecordingSpan) SpanContext() SpanContext { return s.sc }
