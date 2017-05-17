package options

import (
	"fmt"
	"sort"
)

type Options struct {
	opts map[string]Option
}

type Option interface {
	Key() string
	String() string
	Value() interface{}
	Set(value string) error
}

func New() *Options {
	o := Options{}
	o.opts = make(map[string]Option, 0)
	return &o
}

func (o *Options) Add(opt Option) {
	o.opts[opt.Key()] = opt
}

// Keys returns all registered option keys.
func (o *Options) Keys() []string {
	keys := make(sort.StringSlice, 0, len(o.opts))
	for tag := range o.opts {
		keys = append(keys, tag)
	}
	keys.Sort()
	return keys
}

func (o *Options) Get(key string) Option {
	return o.opts[key]
}

func (o *Options) Value(key string) interface{} {
	v := o.Get(key)
	if v == nil {
		return nil
	}
	return v.Value()
}

func (o *Options) StringValue(key string) string {
	val := o.Value(key)
	switch val := val.(type) {
	case string:
		return val
	default:
		panic(fmt.Errorf("Expected string option in StringValue(), got %T", val))
	}
}

func (o *Options) IntValue(key string) int {
	val := o.Value(key)
	switch val := val.(type) {
	case int:
		return val
	default:
		panic(fmt.Errorf("Expected integer option in IntValue(), got %T", val))
	}
}

func (o *Options) BoolValue(key string) bool {
	val := o.Value(key)
	switch val := val.(type) {
	case bool:
		return val
	default:
		panic(fmt.Errorf("Expected boolean option in BoolValue(), got %T", val))
	}
}
