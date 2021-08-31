package kylin

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
}

func (c *Context) Json(data interface{}) (err error) {
	encoder := json.NewEncoder(c.Response)
	if err = encoder.Encode(data); err != nil {
		return err
	}
	return nil
}
