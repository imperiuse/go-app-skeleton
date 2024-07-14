package apierror

// Inspired by this article https://www.linkedin.com/pulse/rfc-7807-error-handling-standard-apis-david-rold%C3%A1n-mart%C3%ADnez
// RFC-7807

// APIError defines the structure for describing API errors, following RFC-7807.
// //nolint: lll, this is usefully len comment line
type APIError struct {
	Type     string `json:"type" example:"reports-service/issues/token-generation-error"` // A URI reference that identifies the problem type.
	Title    string `json:"title" example:"Name of the problem or an error"`              // A short, human-readable summary of the problem type.
	Status   int    `json:"status" example:"500"`                                         // The HTTP status code for this occurrence of the problem.
	Detail   string `json:"detail" example:"Description of the problem"`                  // A human-readable explanation specific to this occurrence.
	Instance string `json:"instance" example:"GET /api/v1/some"`                          // A URI reference that identifies the specific occurrence of the problem.
}
