package samson

type ErrorResponse struct {
	Message string
}

func (er ErrorResponse) Error() string {
	return er.Message
}
