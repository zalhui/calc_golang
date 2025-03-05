package calculation

import (
	"reflect"
	"testing"
)

func TestConvertToRPN(t *testing.T) {
	tests := []struct {
		expression string
		expected   []string
		err        error
	}{
		// Базовые операции
		{"2+2", []string{"2", "2", "+"}, nil},
		{"3-1", []string{"3", "1", "-"}, nil},
		{"4*5", []string{"4", "5", "*"}, nil},
		{"6/2", []string{"6", "2", "/"}, nil},

		// Выражения со скобками
		{"(2+2)*3", []string{"2", "2", "+", "3", "*"}, nil},
		{"2*(3+4)", []string{"2", "3", "4", "+", "*"}, nil},
		{"(5-1)/2", []string{"5", "1", "-", "2", "/"}, nil},

		// Сложные выражения
		{"2+3*4-5/2", []string{"2", "3", "4", "*", "+", "5", "2", "/", "-"}, nil},
		{"(2+3)*(4-1)", []string{"2", "3", "+", "4", "1", "-", "*"}, nil},
		{"1+2+3*4", []string{"1", "2", "+", "3", "4", "*", "+"}, nil},

		// Граничные случаи
		{"42", []string{"42"}, nil},                                       // Одно число
		{"((2+3))", []string{"2", "3", "+"}, nil},                         // Многоуровневые скобки
		{"2*(3*(4+5))", []string{"2", "3", "4", "5", "+", "*", "*"}, nil}, // Вложенные скобки

		// Ошибочные случаи
		{"2++2", nil, ErrValues},     // Два оператора подряд
		{"2+(3*4", nil, ErrBrackets}, // Несбалансированные скобки
		{"(2+3))", nil, ErrBrackets}, // Лишняя закрывающая скобка
		{"2+x", nil, ErrAllowed},     // Недопустимый символ
		{"2 3 +", nil, ErrAllowed},   // Пробелы между числами
		{"", nil, ErrValues},         // Пустое выражение
		{"+", nil, ErrValues},        // Только оператор
		{"(2+3", nil, ErrBrackets},   // Незакрытая скобка
		{"2+3)", nil, ErrBrackets},   // Лишняя закрывающая скобка
		{"2.5+3.7", nil, ErrAllowed}, // Десятичные числа (если не поддерживаются)
	}

	for _, tt := range tests {
		result, err := convertToRPN(tt.expression)
		if !reflect.DeepEqual(result, tt.expected) || err != tt.err {
			t.Errorf("convertToRPN(%q) = %v, %v; want %v, %v", tt.expression, result, err, tt.expected, tt.err)
		}
	}
}
