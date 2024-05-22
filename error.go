package mytrade

import "errors"

var (
	ErrorSymbolNotFound = errors.New("symbol not found")
	ErrorAccountType    = errors.New("account type error")
	ErrorInvalidParam   = errors.New("invalid param")
)
