package mytrade

import "errors"

var (
	ErrorSymbolNotFound = errors.New("symbol not found")
	ErrorAccountType    = errors.New("account type error")
	ErrorInvalidParam   = errors.New("invalid param")
	ErrorNotSupport     = errors.New("not support")
	ErrorOrderNotFound  = errors.New("order not found")
)

func ErrorInvalid(paramName string) error {
	return errors.New("invalid " + paramName)
}
