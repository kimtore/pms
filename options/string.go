package options

import "fmt"

type StringOption struct {
	key   string
	value string
}

func NewStringOption(key string) *StringOption {
	return &StringOption{key: key}
}

func (o *StringOption) Set(value string) error {
	o.value = value
	return nil
}

func (o *StringOption) Key() string {
	return o.key
}

func (o *StringOption) StringValue() string {
	return o.value
}

func (o *StringOption) Value() interface{} {
	return o.value
}

func (o *StringOption) String() string {
	return fmt.Sprintf("%s=%s", o.key, o.value)
}
