package errors

import "errors"

//type DuplicateURLError struct {
//	Shortening dto.Shortening
//	URLError   error
//}

//	func (d *DuplicateURLError) Error() string {
//		return fmt.Sprintf("[%s] %v", d.Shortening, d.URLError)
//	}
var (
	ErrorDuplicateURL = errors.New("duplicate URL")
)
