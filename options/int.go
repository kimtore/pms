package options

import (
	"fmt"
	"strconv"
)

type IntOption struct {
	key   string
	value int
}

func NewIntOption(key string) *IntOption {
	return &IntOption{key: key}
}

func (o *IntOption) Set(value string) error {
	var err error
	o.value, err = strconv.Atoi(value)
	return err
}

func (o *IntOption) Key() string {
	return o.key
}

func (o *IntOption) IntValue() int {
	return o.value
}

func (o *IntOption) Value() interface{} {
	return o.value
}

func (o *IntOption) String() string {
	return fmt.Sprintf("%s=%d", o.key, o.value)
}
