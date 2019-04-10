package jankyError

import "fmt"

//Please use official error handling if it has being released
//https://github.com/golang/proposal/blob/master/design/go2draft-error-handling.md

type TheError struct {
	Code    uint16 //Code cannot less than 0
	Message string
	Detail  interface{} //using Detail to trace stack.
}

var (
	//Error conditions for event bus
	NotDataCode uint16 = 1
	NotData            = "interface is function"
)

//TODO: finish stack trace
func (e *TheError) Error() string {
	return fmt.Sprintf("Error code %v: %s", e.Code, e.Message)
}
