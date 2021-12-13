package kylin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rd-code/kylin/route"
)

type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
	rs       *route.RouterServer
}

func (c *Context) Json(data interface{}) (err error) {
	encoder := json.NewEncoder(c.Response)
	if err = encoder.Encode(data); err != nil {
		return err
	}
	return nil
}

func (c *Context) Parse(data interface{}) (err error) {
	decoder := json.NewDecoder(c.Request.Body)

	if err = decoder.Decode(data); err != nil {
		return
	}
	return
}

func (c *Context) ParamString(key string) string {
	return c.rs.Param(key)
}

func (c *Context) ParamInt(key string) (int, error) {
	v := c.ParamString(key)
	if len(v) == 0 {
		return 0, fmt.Errorf("the param is empty")
	}
	return strconv.Atoi(v)
}

func (c *Context) ParamUint64(key string) (uint64, error) {
	v := c.ParamString(key)
	if len(v) == 0 {
		return 0, fmt.Errorf("the param is empty")
	}
	return strconv.ParseUint(v, 10, 64)
}

func (c *Context) ParamInt64(key string) (int64, error) {
	v := c.ParamString(key)
	if len(v) == 0 {
		return 0, fmt.Errorf("the param is empty")
	}
	return strconv.ParseInt(v, 10, 64)
}
