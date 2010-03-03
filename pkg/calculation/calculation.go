package calculation

import (
	"fmt"
	"strconv"

	//"strconv"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/zalhui/calc_golang/config"
	"github.com/zalhui/calc_golang/internal/orchestrator/models"
)

var priority = map[rune]int{
	'+': 1,
	'-': 1,
	'*': 2,
	'/': 2,
	'(': 0,
}

func ParseExpression(expression string, ExpressionID string) ([]models.Task, error) {
	rpn, err := convertToRPN(expression)

	if err != nil {
		return nil, fmt.Errorf("error converting expression to RPN: %w", err)
	}

	var tasks []models.Task
	var stack []string

	for _, elem := range rpn {
		if isOperator(elem) {
			if len(stack) < 2 {
				return nil, ErrValues
			}
			arg1, arg2 := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			taskID := uuid.NewString()

			value1, err := strconv.ParseFloat(arg1, 64)
			if err != nil {
				return nil, ErrAllowed
			}

			value2, err := strconv.ParseFloat(arg2, 64)
			if err != nil {
				return nil, ErrAllowed
			}
			tasks = append(tasks, models.Task{
				ID:            taskID,
				ExpressionID:  ExpressionID,
				Arg1:          value1,
				Arg2:          value2,
				Operation:     elem,
				OperationTime: getOperationTime(elem),
				Status:        "pending",
			})

			resultPlaceholder := fmt.Sprintf("task_%s_result", taskID)
			stack = append(stack, resultPlaceholder)
		} else {
			stack = append(stack, elem)
		}
	}
	return tasks, nil
}

func isOperator(r string) bool {
	return r == "+" || r == "-" || r == "*" || r == "/"
}

func getOperationTime(r string) time.Duration {
	cfg := config.LoadConfig()
	switch r {
	case "+":
		return cfg.TimeAddition
	case "-":
		return cfg.TimeSubtraction
	case "*":
		return cfg.TimeMultiplication
	case "/":
		return cfg.TimeDivision
	}
	return 0
}

/*func Calc(expression string) (float64, error) {
	rpn, err := convertToRPN(expression)

	if err != nil {
		return 0, err
	}

	return calculateRPN(rpn)
}*/

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
				return nil, ErrAllowed
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

/*func calculateRPN(rpn []string) (float64, error) {
	var stack []float64

	for _, elem := range rpn {
		switch elem {
		case "+", "-", "*", "/":
			if len(stack) < 2 {
				return 0, ErrValues
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
				return 0, ErrAllowed
			}
			stack = append(stack, value)
		}
	}

	if len(stack) != 1 {
		return 0, ErrValues
	}

	return stack[0], nil
}*/
