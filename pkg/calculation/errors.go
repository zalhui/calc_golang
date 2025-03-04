package calculation

import "errors"

var (
	ErrBrackets       = errors.New("expression is not valid. number of brackets doesn't match")
	ErrValues         = errors.New("expression is not valid. not enough values")
	ErrDivisionByZero = errors.New("expression is not valid. division by zero")
	ErrAllowed        = errors.New("expression is not valid. only numbers and ( ) + - * / allowed")
)
