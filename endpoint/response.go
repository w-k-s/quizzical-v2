package endpoint

type Response struct{
	Data interface{} `json:"data"`
}

/*
An Error response is supposed to be serialized to the following form:
{"error":<object|string>
Error Wrapper is used to ensure the root error field is present.
*/
type ErrorWrapper struct{
	OriginalError error `json:"error"`
}

func (e ErrorWrapper) Error() string{
	return e.OriginalError.Error()
}