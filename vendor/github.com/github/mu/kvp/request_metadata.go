package kvp

import "github.com/github/mu/ctxkey"

// Tags is an alias of map[string]string, a type for tags associated with a metric.
type Tags map[string]string

// Merge creates a merged set of Tags. Tags passed in will overwite current
// tags.
func (t Tags) Merge(tags Tags) Tags {
	merged := make(Tags)
	for k, v := range t {
		merged[k] = v
	}
	for k, v := range tags {
		merged[k] = v
	}

	return merged
}

// RMDContextKey is a key that can be used for adding and retrieving a
// RequestMeatadata from a Context.
var RMDContextKey = ctxkey.New("RequestMetadata")

// RequestMetadata holds a set of kvp's to add to log output and a set of
// Tags to add to metrics.
type RequestMetadata struct {
	logFields *KVP
	statTags  Tags
	skipLog   bool
}

// NewRequestMetadata initializes a new RequestMetadata
func NewRequestMetadata() *RequestMetadata {
	return &RequestMetadata{
		logFields: New(),
		statTags:  Tags{},
	}
}

// Copy returns a new RequestMetadata copying the fields and tags.
func (m *RequestMetadata) Copy() *RequestMetadata {
	return &RequestMetadata{
		logFields: New(m.logFields.Fields()...),
		statTags:  m.statTags,
		skipLog:   m.skipLog,
	}
}

// LogWith adds the fields to the set of logging kvps.
func (m *RequestMetadata) LogWith(fields ...Field) {
	m.logFields.Add(fields...)
}

// TagStatsWith adds the tags to the set of metrics tags.
func (m *RequestMetadata) TagStatsWith(tags Tags) {
	m.statTags = m.statTags.Merge(tags)
}

// LogFields returns the kvp fields for logging
func (m *RequestMetadata) LogFields() []Field {
	return m.logFields.Fields()
}

// StatTags returns the metrics tags
func (m *RequestMetadata) StatTags() Tags {
	return m.statTags
}

// SkipLogging turns off HTTP request logging for this request.
func (m *RequestMetadata) SkipLogging() {
	m.skipLog = true
}

// ShouldSkipLogging returns true if SkipLogging was called for this request's context.
func (m *RequestMetadata) ShouldSkipLogging() bool {
	return m.skipLog
}
