package errors

type CustomError struct {
	StatusCode int
	Msg        string
}

func (e *CustomError) Error() string {
	return e.Msg
}

func (e *CustomError) BadRequest() *CustomError {
	e.StatusCode = 400
	return e

}

func (e *CustomError) Unauthorized() *CustomError {
	e.StatusCode = 401
	return e
}

func (e *CustomError) InternalServerError() *CustomError {
	e.StatusCode = 500
	return e
}

func Cause() *CustomError {
	return &CustomError{}
}

func (e *CustomError) New(msg string) *CustomError {
	e.Msg = msg
	return e
}
