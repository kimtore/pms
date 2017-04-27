package options

import "fmt"

type StringOption struct {
	key   string
	value string
}

func NewStringOption(key string, value string) *StringOption {
	o := &StringOption{key: key}
	err := o.Set(value)
	if err != nil {
		panic(err)
	}
	return o
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
