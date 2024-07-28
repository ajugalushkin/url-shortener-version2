package errors

import "errors"

// ErrorDuplicateURL Ошибка дублирования
var (
	ErrorDuplicateURL = errors.New("duplicate URL")
)
