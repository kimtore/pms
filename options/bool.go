package options

import "strconv"

type BoolOption struct {
	key   string
	value bool
}

func NewBoolOption(key string) *BoolOption {
	return &BoolOption{key: key}
}

func (o *BoolOption) Set(value string) error {
	var err error
	o.value, err = strconv.ParseBool(value)
	return err
}

func (o *BoolOption) SetBool(value bool) {
	o.value = value
}

func (o *BoolOption) Key() string {
	return o.key
}

func (o *BoolOption) BoolValue() bool {
	return o.value
}

func (o *BoolOption) Value() interface{} {
	return o.value
}

func (o *BoolOption) String() string {
	t := o.Key()
	if !o.value {
		t = "no" + t
	}
	return t
}

func (o *BoolOption) StringValue() string {
	return o.String()
}
