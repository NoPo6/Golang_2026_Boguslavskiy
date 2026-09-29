package calculator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalculate_Success(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		expected   float64
	}{
		{"addition", "1+0", 1},
		{"subtraction", "101-3", 98},
		{"multiplication", "22*3", 66},
		{"division", "10/2", 5},
		{"division fractional", "7/2", 3.5},
		{"negative result", "3-5", -2},
		{"precedence", "1+2*3", 7},
		{"brackets", "(1+2)*3", 9},
		{"nested brackets", "((1+2)*3)-4", 5},
		{"unary minus", "-5", -5},
		{"unary minus with brackets", "-(3+2)", -5},
		{"multiple unary minus (even)", "-(-(-(-5)))", 5},
		{"multiple unary minus (odd)", "-(-(-5))", -5},
		{"implicit multiplication brackets", "(9+3)(7+2)", 108},
		{"implicit multiplication number and bracket", "2(3+4)", 14},
		{"implicit multiplication bracket and number", "(2+3)4", 20},
		{"example from task", "(1+2)-3", 0},
		{"example from task two", "(1+2)*3", 9},
		{"with spaces", "(1 + 2) - 3", 0},
		{"float numbers", "1.5+2.5", 4},
		{"complex", "2*(3+4)-10/5", 12},
		{"division and multiplication chain", "100/5/2*3", 30},
		{"multiplication with unary minus", "5*(-2)", -10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Calculate(tt.expression)
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculate_Failure(t *testing.T) {
	tests := []struct {
		name       string
		expression string
	}{
		{"empty expression", ""},
		{"only spaces", "   "},
		{"unmatched open bracket", "(1+2"},
		{"unmatched close bracket", "1+2)"},
		{"division by zero", "1/0"},
		{"invalid char", "1+a"},
		{"trailing operator", "1+"},
		{"leading operator", "*1"},
		{"double plus", "1++2"},
		{"double unary minus no brackets", "--5"},
		{"triple unary minus no brackets", "---5"},
		{"binary and unary minus no brackets", "5--5"},
		{"two dots in number", "1.1.1"},
		{"many dots in number", "1.......5"},
		{"double division", "5//5"},
		{"double multiplication", "5**5"},
		{"operator at start of brackets division", "2/(/2)"},
		{"operator at start of brackets multiplication", "2/(*2)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Calculate(tt.expression)
			require.Error(t, err)
		})
	}
}
