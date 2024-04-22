package api

type Error struct {
	// The HTTP status code of the error.
	Status int `json:"status"`
	// A machine-readable code indicating the type of error.
	Code string `json:"code"`
	// A summary of the error.
	Title string `json:"title,omitempty"`
	// A detailed explanation of the error.
	Detail string    `json:"detail,omitempty"`
	Source *struct { // The query parameter that produced the error.
		Parameter string `json:"parameter,omitempty"`
		// A JSON pointer that indicates the location in the request entity where the error originates.
		Pointer string `json:"pointer,omitempty"`
	} `json:"source,omitempty"`
}

func (e *Error) Error() string {
	return e.Title + ": " + e.Detail
}

type ErrorsResponse struct {
	Errors []Error `json:"errors"`
}

type ForbiddenResponse struct {
	Errors []Error
}

func NewInternalServerError() Error {
	return Error{
		Status: 500,
		Code:   "Internal Server Error",
		Title:  "",
		Detail: "",
	}
}
