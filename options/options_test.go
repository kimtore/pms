package options_test

import (
	"testing"

	"github.com/ambientsound/pms/options"
	"github.com/stretchr/testify/require"
)

func TestStringOption(t *testing.T) {
	opt := options.NewStringOption("foo")
	err := opt.Set("bar")
	require.Nil(t, err)
	if opt.Key() != "foo" {
		t.Fatalf("String option key is incorrect!")
	}
	val := opt.Value()
	switch val := val.(type) {
	case string:
		if val != "bar" {
			t.Fatalf("String option value is incorrect!")
		}
	default:
		t.Fatalf("String option value is of wrong type %T!", val)
	}
}

func TestIntOption(t *testing.T) {
	opt := options.NewIntOption("foo")
	err := opt.Set("3984")
	require.Nil(t, err)
	if opt.Key() != "foo" {
		t.Fatalf("Int option key is incorrect!")
	}
	val := opt.Value()
	switch val := val.(type) {
	case int:
		if val != 3984 {
			t.Fatalf("Int option value is incorrect!")
		}
	default:
		t.Fatalf("Int option value is of wrong type %T!", val)
	}
}

func TestBoolOption(t *testing.T) {
	opt := options.NewBoolOption("foo")
	err := opt.Set("true")
	require.Nil(t, err)
	if opt.Key() != "foo" {
		t.Fatalf("Bool option key is incorrect!")
	}
	val := opt.Value()
	switch val := val.(type) {
	case bool:
		if !val {
			t.Fatalf("Bool option value is incorrect!")
		}
	default:
		t.Fatalf("Bool option value is of wrong type %T!", val)
	}
}
