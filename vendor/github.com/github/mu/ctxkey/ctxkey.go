package ctxkey

// The basic idea here comes from https://github.com/cstockton/pkg/tree/master/ctxkey

// Key holds a globally unique named value for use with contexts.
type Key interface {
	Name() string
}

// New creates a new Key that will not cause allocations when used.
func New(name string) Key {
	return &key{name}
}

type key struct {
	name string
}

func (k *key) Name() string {
	return k.name
}
