package main

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	testCase := struct {
		values   []string
		expected map[string]int
	}{
		values: []string{
			"foo",
			"bar",
			"foo",
			"foo",
			"fizz",
			"bar",
		},
		expected: map[string]int{
			"foo":  3,
			"bar":  2,
			"fizz": 1,
		},
	}

	c := NewCounter()
	for _, s := range testCase.values {
		_ = c.Add(s)
	}

	for k, v := range testCase.expected {
		count, ok := c.v[k]
		if !ok {
			t.Errorf("%v not found\n", k)
		}

		if count != v {
            t.Errorf("Expected %v =  %v, found %v\n", k, v, count)
		}

		fmt.Printf("Expected %v = %v, found %v\n", k, v, count)
	}
}
