package http_response

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Errors     interface{} `json:"errors,omitempty"`
	LastPage   int         `json:"last_page,omitempty"`
}

type ResponseHealthCheck struct {
	Message string `json:"message"`
}
