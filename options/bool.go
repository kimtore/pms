package options

import "strconv"

type BoolOption struct {
	key   string
	value bool
}

func NewBoolOption(key string, value bool) (o *BoolOption, err error) {
	o = &BoolOption{}
	o.key = key
	o.SetBool(value)
	return
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

func (o *BoolOption) Text() string {
	t := o.Key()
	if !o.value {
		t = "no" + t
	}
	return t
}
