package error

type code int

const (
	InvalidJSON code = iota
)


type HTTPError struct {
	ErrorMessage string 	`json:"message,omitempty"`
	ErrorCode code			`json:"code"`
}
