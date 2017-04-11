package options_test

import (
	"testing"

	"github.com/ambientsound/pms/options"
)

func TestOptions(t *testing.T) {
	var opt options.Option
	opts := options.New()
	opt, _ = options.NewStringOption("foo", "bar")
	opts.Add(opt)
	if opts.StringValue("foo") != "bar" {
		t.Fatalf("String value was not correctly added to options!")
	}
}

func TestStringOption(t *testing.T) {
	opt, err := options.NewStringOption("foo", "bar")
	if err != nil {
		t.Fatalf("NewStringOption() failed initialization: %s", err)
	}
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
	opt, err := options.NewIntOption("foo", "3984")
	if err != nil {
		t.Fatalf("NewIntOption() failed initialization: %s", err)
	}
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
	opt, err := options.NewBoolOption("foo", true)
	if err != nil {
		t.Fatalf("NewBoolOption() failed initialization: %s", err)
	}
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
