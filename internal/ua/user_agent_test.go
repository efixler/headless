package ua

import (
	"testing"
)

func TestArg(t *testing.T) {
	type data struct {
		name     string
		in       string
		expected string
	}
	tests := []data{
		{"Firefox88", ":firefox:", Firefox88},
		{"Safari537", ":safari:", Safari537},
		{"Custom", "custom", "custom"},
	}
	for _, test := range tests {
		a := &Arg{}
		err := a.UnmarshalText([]byte(test.in))
		if err != nil {
			t.Errorf("[%s] Error unmarshalling %s: %s", test.name, test.in, err)
		}
		if a.String() != test.expected {
			t.Errorf("[%s] Expected %s, got %s", test.name, test.expected, a.String())
		}
	}
}
