package utils

//自定义了error
type MyError struct {
	Code    int
	Message string
}

func (this *MyError) Error() string {
	return this.Message
}

//自定义error
func NewMyError(code int, msg string) error {
	return &MyError{Code: code, Message: msg}
}
