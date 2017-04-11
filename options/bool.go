package options

type BoolOption struct {
	key   string
	value bool
}

func NewBoolOption(key string, value bool) (o *BoolOption, err error) {
	o = &BoolOption{}
	o.key = key
	err = o.Set(value)
	return
}

func (o *BoolOption) Set(value bool) error {
	o.value = value
	return nil
}

func (o *BoolOption) Key() string {
	return o.key
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
