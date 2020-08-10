package apiErrors

type NotFoundError struct {
	S string
}

func (e *NotFoundError) Error() string {
	return e.S
}
