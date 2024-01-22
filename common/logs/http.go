package logs

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewGinStructuredLogger(logger *logrus.Logger) gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: newStructuredLogger(logger),
	})
}

func newStructuredLogger(logger *logrus.Logger) gin.LogFormatter {
	return func(param gin.LogFormatterParams) string {
		var logFields = make(map[string]interface{})

		logFields["http_method"] = param.Method
		logFields["uri"] = param.Path
		logFields["status"] = param.StatusCode
		logFields["latency"] = param.Latency
		logFields["client_ip"] = param.ClientIP
		logFields["error"] = param.ErrorMessage

		if len(param.ErrorMessage) > 0 {
			logFields["errors"] = param.ErrorMessage
		}

		entry := logger.WithFields(logFields)

		msg := fmt.Sprintf("[%s] %s %d %s %s",
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.StatusCode,
			param.Path,
			param.Latency,
		)

		if param.ErrorMessage != "" {
			msg = msg + " " + param.ErrorMessage
		}

		if len(param.ErrorMessage) > 0 {
			msg = msg + " " + param.ErrorMessage
		}

		entry.Info(msg)
		return ""
	}
}
