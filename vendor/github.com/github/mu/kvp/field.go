package kvp

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

// String creates a Field that stores a string value.
func String(key, val string) Field {
	return Field{t: 's', key: key, str: val}
}

// Error creates a Field that stores an error message.
func Error(val error) Field {
	return Field{t: 's', key: "error", str: val.Error()}
}

// Int creates a Field that stores an int value.
func Int(key string, val int) Field {
	return Field{t: 'i', key: key, int: int64(val)}
}

// Int64 creates a Field that stores an int value.
func Int64(key string, val int64) Field {
	return Field{t: 'i', key: key, int: val}
}

// Float creates a Field stat stores a float64 value.
func Float(key string, val float64) Field {
	return Field{t: 'f', key: key, int: int64(math.Float64bits(val))}
}

// Bool creates a Field that stores a bool value.
func Bool(key string, val bool) Field {
	return Field{t: 'b', key: key, boolean: val}
}

// Duration creates a Field that stores a time.Duration value.
func Duration(key string, val time.Duration) Field {
	return Field{t: 'd', key: key, int: int64(val)}
}

// Time creates a Field that stores a time.Time value.
func Time(key string, time time.Time) Field {
	return Field{t: 't', key: key, int: time.UnixNano()}
}

// Any creates a Field that can story any value.
func Any(key string, any interface{}) Field {
	return Field{t: 'a', key: key, any: any}
}

// Field is a generic type that stores a key-value pair entry to be associated
// with the KVP.
type Field struct {
	key     string
	str     string
	int     int64
	boolean bool
	t       byte
	any     interface{}
}

// Key returns the key for the field.
func (f *Field) Key() string {
	return f.key
}

// String returns the value stored in this field as a string, if it can be
// converted. This should be use sparingly and is only here until the haystack
// reporter takes a real context instead of a map[string]string.
func (f *Field) String() string {
	switch f.t {
	case 'i':
		return strconv.FormatInt(f.int, 10)
	case 'f':
		return strconv.FormatFloat(f.asFloat(), 'E', -1, 64)
	case 'b':
		return fmt.Sprintf("%t", f.boolean)
	case 'd':
		return f.asDuration().String()
	case 't':
		return f.asTime().Format("2006-01-02 15:04:05.999999999 -0700 MST")
	case 'a':
		return fmtAny(f.any)
	default:
		return f.str
	}
}

// Value returns the raw value
func (f *Field) Value() interface{} {
	switch f.t {
	case 'i':
		return f.int
	case 'f':
		return f.asFloat()
	case 'b':
		return f.boolean
	case 'd':
		return f.asDuration()
	case 't':
		return f.asTime()
	case 'a':
		return f.any
	default:
		return f.str
	}
}

// asFloat returns the float64 value stored in this Field, if the Field is a
// float.
func (f *Field) asFloat() float64 {
	return math.Float64frombits(uint64(f.int))
}

// asDuration returns the time.Duration value stored in this field, if the
// field is a duration.
func (f *Field) asDuration() time.Duration {
	return time.Duration(f.int)
}

// asTime returns the time.Time stored in this field, if the field is a time.
func (f *Field) asTime() time.Time {
	return time.Unix(f.int/1e9, f.int%1e9)
}

func fmtAny(any interface{}) string {
	switch str := any.(type) {
	case fmt.Stringer:
		return str.String()
	case fmt.GoStringer:
		return str.GoString()
	default:
		return fmt.Sprint(str)
	}
}
