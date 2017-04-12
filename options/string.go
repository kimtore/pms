package options

type StringOption struct {
	key   string
	value string
}

func NewStringOption(key string, value string) (o *StringOption, err error) {
	o = &StringOption{}
	o.key = key
	err = o.Set(value)
	return
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

func (o *StringOption) Text() string {
	return o.value
}
