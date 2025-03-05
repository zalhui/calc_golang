package worker

import (
	"testing"

	"github.com/zalhui/calc_golang/pkg/calculation"
)

func TestPerformOperation(t *testing.T) {
	tests := []struct {
		arg1      float64
		arg2      float64
		operation string
		expected  float64
		err       error
	}{
		{2, 3, "+", 5, nil},
		{5, 3, "-", 2, nil},
		{4, 2, "*", 8, nil},
		{6, 3, "/", 2, nil},
		{6, 0, "/", 0, calculation.ErrDivisionByZero},
	}

	for _, tt := range tests {
		result, err := performOperation(tt.arg1, tt.arg2, tt.operation)
		if result != tt.expected || err != tt.err {
			t.Errorf("performOperation(%v, %v, %q) = %v, %v; want %v, %v", tt.arg1, tt.arg2, tt.operation, result, err, tt.expected, tt.err)
		}
	}
}
