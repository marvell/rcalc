package main

import (
	"reflect"
	"testing"
	"time"
)

func TestParseDistance(t *testing.T) {
	tests := []struct {
		in  string
		out float64
	}{
		{"10km", 10},
		{"1km", 1},
		{"0.1km", 0.1},
		{"", 0},
	}

	for _, test := range tests {
		got, _ := ParseDistance(test.in)
		if got != test.out {
			t.Errorf("ParseDistance(%s) returns %f, but wants %f", test.in, got, test.out)
		}
	}
}

func TestParseExpression(t *testing.T) {
	tests := []struct {
		in  string
		out TDP
	}{
		{
			"10km*4m30s",
			TDP{
				Distance: 10,
				Pace:     4*time.Minute + 30*time.Second,
				Time:     45 * time.Minute,
			},
		},
		{
			"45m/4m30s",
			TDP{
				Distance: 10,
				Pace:     4*time.Minute + 30*time.Second,
				Time:     45 * time.Minute,
			},
		},
		{
			"45m/10km",
			TDP{
				Distance: 10,
				Pace:     4*time.Minute + 30*time.Second,
				Time:     45 * time.Minute,
			},
		},
		{
			"3m30s/0.6km",
			TDP{
				Distance: 0.6,
				Pace:     5*time.Minute + 50*time.Second,
				Time:     3*time.Minute + 30*time.Second,
			},
		},
	}

	for _, test := range tests {
		got, err := ParseExpression(test.in)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(*got, test.out) {
			t.Errorf("ParseExpression(%s) returns %+v, but wants %+v", test.in, *got, test.out)
		}
	}
}
