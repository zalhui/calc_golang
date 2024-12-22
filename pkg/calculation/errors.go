package calculation

import "errors"

var (
	//ErrInvalidExpression = errors.New("Expression is not valid. ")
	ErrBrackets       = errors.New("Expression is not valid. Number of brackets doesn't match")
	ErrValues         = errors.New("Expression is not valid. Not enough values")
	ErrDivisionByZero = errors.New("Expression is not valid. Division by zero")
	ErrAllowed        = errors.New("Expression is not valid. Only numbers and ( ) + - * / allowed")
)
