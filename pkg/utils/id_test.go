package utils

import (
	"fmt"
	"testing"
)

func TestNewNumericID(t *testing.T) {
	type testCase struct {
		Name string

		ExpectedInt int
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			for i := 0; i < 10; i++ {
				fmt.Println(NewNumericID())
			}
		})
	}

	validate(t, &testCase{})
}
