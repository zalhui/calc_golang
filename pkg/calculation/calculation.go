package calculation

import (
	"fmt"

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
		if isOperator(rune(elem[0])) {
			if len(stack) < 2 {
				return nil, ErrValues
			}
			arg1, arg2 := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			taskID := uuid.NewString()

			tasks = append(tasks, models.Task{
				ID:            taskID,
				ExpressionID:  ExpressionID,
				Arg1:          arg1,
				Arg2:          arg2,
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
	fmt.Printf("Tasks: %+v\n", tasks)
	return tasks, nil
}

func isOperator(r rune) bool {
	return r == '+' || r == '-' || r == '*' || r == '/'
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

	fmt.Println(rpn) // Для отладки выводим RPN
	return rpn, nil
}
