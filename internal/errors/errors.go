package errors

import (
	"fmt"

	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
)

type DuplicateURLError struct {
	Shortening dto.Shortening
	URLError   error
}

func (d *DuplicateURLError) Error() string {
	return fmt.Sprintf("[%s] %v", d.Shortening, d.URLError)
}
