package dtos

type ClientErrorType = string // @name ErrorType
const (
	ValidationErrorType ClientErrorType = "validation"
	DatabaseErrorType   ClientErrorType = "database"
)

// ClientError is an error that can be returned to the client.
type ClientError struct {
	Code          int           `json:"code" example:"400" extensions:"x-order=1"`                                                                       // Similar to the http status code
	Type          string        `json:"type,omitempty" enums:"validation,database,entity" example:"validation" extensions:"x-order=2"`                   // The type of error
	Message       string        `json:"message" example:"Invalid request Body" extensions:"x-order=3"`                                                   // A human-readable error message
	Details       []interface{} `json:"details,omitempty" swaggertype:"array,object" extensions:"x-order=4"`                                             // Additional details about the error, omitted if empty
	CorrelationId string        `json:"correlation_id,omitempty" example:"01234567890123456789012345678900" swaggertype:"string" extensions:"x-order=5"` // The error id from context traceId
} // @name Error

// used to implement the error interface
func (c ClientError) Error() string {
	return c.Message
}
