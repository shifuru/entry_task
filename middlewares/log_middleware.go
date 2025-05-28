package middlewares

import (
	"bytes"
	"encoding/json"
	"entry_task/common/log"
	"entry_task/common/utils"
	"github.com/gin-gonic/gin"
	"io"
	"sync"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

var (
	//缓冲区池
	bufferPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
)

type AccessParam struct {
	StartTime time.Time       `json:"start_time"`
	Status    int             `json:"status"`
	UserId    uint            `json:"user_id"`
	Username  string          `json:"username"`
	UserGroup uint            `json:"user_group"`
	Ip        string          `json:"ip"`
	Method    string          `json:"method"`
	Path      string          `json:"path"`
	Request   json.RawMessage `json:"request"`
	Response  json.RawMessage `json:"response"`
}

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		buffer := bufferPool.Get().(*bytes.Buffer)
		buffer.Reset()

		blw := &bodyLogWriter{body: buffer, ResponseWriter: c.Writer}
		c.Writer = blw
		startTime := time.Now()

		request, err := io.ReadAll(c.Request.Body)
		if nil == err {
			_ = c.Request.Body.Close()
			c.Request.Body = io.NopCloser(bytes.NewBuffer(request))
		}

		var buf bytes.Buffer
		err = json.Compact(&buf, request)
		if nil == err {
			request = buf.Bytes()
		}
		if len(request) == 0 {
			request = []byte("{}")
		}

		c.Next()

		accessParam := AccessParam{
			StartTime: startTime,
			Status:    blw.Status(),
			UserId:    c.GetUint("user_id"),
			Username:  c.GetString("username"),
			UserGroup: c.GetUint("user_group"),
			Ip:        c.ClientIP(),
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Request:   request,
			Response:  json.RawMessage(blw.body.String()),
		}

		bufferPool.Put(buffer)

		jsonStr, err := json.Marshal(accessParam)
		if nil != err {
			log.Access(utils.GetRequestId(c), "request: %+v", accessParam)
		} else {
			log.Access(utils.GetRequestId(c), "request: %s", jsonStr)
		}

	}
}
