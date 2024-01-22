package http_response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorResponseHandler(c *gin.Context, statusCode int, err error, data ...interface{}) {
	if err.Error() == "circuit breaker is open" {
		statusCode = 503
	}
	var HttpStatus int
	switch statusCode {
	case 500:
		HttpStatus = http.StatusInternalServerError
	case 503:
		HttpStatus = http.StatusServiceUnavailable
	case 400:
		HttpStatus = http.StatusBadRequest
	case 401:
		HttpStatus = http.StatusUnauthorized
	case 404:
		HttpStatus = http.StatusNotFound
	case 422:
		HttpStatus = http.StatusUnprocessableEntity
	default:
		HttpStatus = http.StatusNotImplemented
	}
	c.AbortWithStatusJSON(HttpStatus, Response{
		StatusCode: statusCode,
		Message:    http.StatusText(HttpStatus),
		Data:       data,
		Errors:     err.Error(),
	})
}

func HttpResponseHandler(c *gin.Context, statusCode int, message string, data interface{}) {
	var HttpStatus int
	switch statusCode {
	case 200:
		HttpStatus = http.StatusOK
	case 201:
		HttpStatus = http.StatusCreated
	case 202:
		HttpStatus = http.StatusAccepted
	default:
		HttpStatus = http.StatusOK
	}
	c.JSON(HttpStatus, Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Errors:     nil,
	})
}

func HttpResponseHandlerWithTotalPage(c *gin.Context, statusCode int, message string, data interface{}, totalPage int) {
	var HttpStatus int
	switch statusCode {
	case 200:
		HttpStatus = http.StatusOK
	case 201:
		HttpStatus = http.StatusCreated
	case 202:
		HttpStatus = http.StatusAccepted
	default:
		HttpStatus = http.StatusOK
	}
	c.JSON(HttpStatus, Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Errors:     nil,
		LastPage:   totalPage,
	})
}
