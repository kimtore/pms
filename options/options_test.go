package options_test

import (
	"testing"

	"github.com/ambientsound/pms/options"
)

func TestOptions(t *testing.T) {
	opts := options.New()
	opts.Add(options.NewStringOption("foo", "bar"))
	if opts.StringValue("foo") != "bar" {
		t.Fatalf("String value was not correctly added to options!")
	}
}

func TestStringOption(t *testing.T) {
	opt := options.NewStringOption("foo", "bar")
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
	opt := options.NewIntOption("foo", 3984)
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
	opt := options.NewBoolOption("foo", true)
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
