package calculation

import (
	"strconv"
	"unicode"
)

var priority = map[rune]int{
	'+': 1,
	'-': 1,
	'*': 2,
	'/': 2,
	'(': 0,
}

func Calc(expression string) (float64, error) {
	rpn, err := convertToRPN(expression)

	if err != nil {
		return 0, err
	}

	return calculateRPN(rpn)
}

// RPN - reverse polish notation
func convertToRPN(expression string) ([]string, error) {
	var rpn []string
	var operators []rune

	// op - operator
	pushOperator := func(op rune) {
		for len(operators) > 0 && priority[operators[len(operators)-1]] >= priority[op] {
			rpn = append(rpn, string(operators[len(operators)-1]))
			operators = operators[:len(operators)-1]
		}
		operators = append(operators, op)
	}

	i := 0
	for i < len(expression) {
		char := rune(expression[i])

		if unicode.IsDigit(char) || char == '.' {
			j := i
			for i < len(expression) && (unicode.IsDigit(rune(expression[i])) || rune(expression[i]) == '.') {
				i++
			}
			rpn = append(rpn, expression[j:i])
			continue
		}

		switch char {
		case '+', '-', '/', '*':
			pushOperator(char)
		case '(':
			operators = append(operators, char)
		case ')':
			for len(operators) > 0 && operators[len(operators)-1] != '(' {
				rpn = append(rpn, string(operators[len(operators)-1]))
				operators = operators[:len(operators)-1]
			}
			if len(operators) == 0 {
				return nil, ErrBrackets
			}
			operators = operators[:len(operators)-1] // удаляем '('
		default:
			if !unicode.IsSpace(char) {
				return nil, ErrInvalidExpression
			}
		}
		i++
	}

	for len(operators) > 0 {
		if operators[len(operators)-1] == '(' {
			return nil, ErrBrackets
		}
		rpn = append(rpn, string(operators[len(operators)-1]))
		operators = operators[:len(operators)-1]
	}

	//fmt.Println(rpn) // Для отладки выводим RPN
	return rpn, nil
}

func calculateRPN(rpn []string) (float64, error) {
	var stack []float64

	for _, elem := range rpn {
		switch elem {
		case "+", "-", "*", "/":
			if len(stack) < 2 {
				return 0, ErrInvalidExpression
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			var result float64
			switch elem {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, ErrDivisionByZero
				}
				result = a / b
			}
			stack = append(stack, result)
		default:
			// convert string to float64
			value, err := strconv.ParseFloat(elem, 64)
			if err != nil {
				return 0, ErrInvalidExpression
			}
			stack = append(stack, value)
		}
	}

	if len(stack) != 1 {
		return 0, ErrInvalidExpression
	}

	return stack[0], nil
}
