package main

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestCheckTargetYmd(t *testing.T) {
	tests := []struct {
		in  string
		out error
	}{
		{"20180610", nil},
		{"20181232", errors.New("targetYmd is bad format")},
	}

	for _, v := range tests {
		err := checkTargetYmd(v.in)
		if err != nil {
			if !reflect.DeepEqual(err.Error(), v.out.Error()) {
				t.Errorf("input: %v\n, get: %#v\n, want: %#v\n", v.in, err.Error(), v.out.Error())
			}
		} else {
			if err != v.out {
				t.Errorf("input: %v\n, get: %v\n, want: %v\n", v.in, err, v.out)
			}
		}
	}
}
