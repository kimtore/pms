package options

import "strconv"

type IntOption struct {
	key   string
	value int
}

func NewIntOption(key string, value string) (o *IntOption, err error) {
	o = &IntOption{}
	o.key = key
	err = o.Set(value)
	return
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

func (o *IntOption) Text() string {
	return string(o.value)
}
