package model

import (
	"fmt"
)

type On interface {
	fmt.Stringer
	isOn()
}

type OnEvent struct {
	Event string
}

type OnSchedule struct {
	Expression string
}

type OnInvalid struct {
	Raw string
}

func (o *OnEvent) isOn()    {}
func (o *OnSchedule) isOn() {}
func (o *OnInvalid) isOn()  {}

func (o *OnEvent) String() string {
	return o.Event
}

func (o *OnSchedule) String() string {
	return fmt.Sprintf("schedule(%s)", o.Expression)
}

func (o *OnInvalid) String() string {
	return o.Raw
}
