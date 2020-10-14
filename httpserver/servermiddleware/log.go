package servermiddleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"xj_web_server/util"
	"time"
)

type bufferedWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bufferedWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}
func BaseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &bufferedWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		// Process request
		c.Next()
		// Log only when path is not being skipped
		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		//var statusColor, methodColor, resetColor string
		//if isTerm {
		//	statusColor = colorForStatus(statusCode)
		//	methodColor = colorForMethod(method)
		//	resetColor = reset
		//}
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()
		if raw != "" {
			path = path + "?" + raw
		}
		req, _ := c.GetRawData()
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(req)) // 关键点
		if len(req) == 0 {
			req = []byte(c.GetString("req"))
		}
		if statusCode == http.StatusInternalServerError {
			util.Logger.Errorf("[QP WEB SERVER] |%3d|%s|%4s|%s|%s|%s|Req:%s|Resp:%s",
				statusCode,
				latency,
				clientIP,
				method,
				path,
				comment,
				string(req),
				string(blw.body.Bytes(),
				),
			)
		} else {
			util.Logger.Infof("[QP WEB SERVER] |%3d|%s|%4s|%s|%s|%s|Req:%s|Resp:%s",
				statusCode,
				latency,
				clientIP,
				method,
				path,
				comment,
				string(req),
				string(blw.body.Bytes()),
			)
		}

	}
}
