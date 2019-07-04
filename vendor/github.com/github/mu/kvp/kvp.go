// Package kvp provides a rich context for loggers and reporters
package kvp

// KVP provides a context that can be used by loggers and reporters.
type KVP struct {
	kvps []Field
}

// New builds a new KVP with the given fields.
func New(fields ...Field) *KVP {
	kvps := make([]Field, 0, len(fields))
	kvps = append(kvps, fields...)
	return &KVP{
		kvps: kvps,
	}
}

// With creates a KVP with the given fields.
func (k *KVP) With(fields ...Field) *KVP {
	kvps := make([]Field, 0, len(k.kvps)+len(fields))
	kvps = append(kvps, k.kvps...)
	kvps = append(kvps, fields...)

	return &KVP{
		kvps: kvps,
	}
}

// Add adds fields to a KVP
func (k *KVP) Add(fields ...Field) {
	k.kvps = append(k.kvps, fields...)
}

// Fields returns the fields.
func (k *KVP) Fields() []Field {
	return k.kvps
}

// Field attempts to find a Field with the matching key name. If no field is
// found, nil is returned.
func (k *KVP) Field(key string) *Field {
	for _, f := range k.kvps {
		if f.Key() == key {
			return &f //nolint: scopelint
		}
	}
	return nil
}
