package kylin

import (
	"github.com/rd-code/kylin/route"
	"net/http"
)

type Engine struct {
	router route.Router
}

func NewEngine() *Engine {
	return &Engine{
		router: route.NewRouter(),
	}
}

func (e *Engine) Get(path string, f func(ctx *Context)) {
	e.router.Get(path, f)
}

func (e *Engine) Post(path string, f func(ctx *Context)) {
	e.router.Post(path, f)
}

func (e *Engine) Delete(path string, f func(ctx *Context)) {
	e.router.Delete(path, f)
}

func (e *Engine) Put(path string, f func(ctx *Context)) {
	e.router.Delete(path, f)
}

func (e *Engine) Handle(path string, f func(ctx *Context)) {
	e.router.Handle(path, f)
}

func (e *Engine) Group(path string) *Engine {
	res := NewEngine()
	res.router = e.router.Group(path)
	return res
}

func (e *Engine) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	handler := e.router.GetHandler(request.Method, request.URL.Path)
	if handler == nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	f := handler.(func(ctx *Context))
	ctx := &Context{
		Request:  request,
		Response: response,
	}
	f(ctx)
}

func (e *Engine) Listen(addr string) {
	http.ListenAndServe(addr, e)
}
