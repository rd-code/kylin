package kylin

import (
	"net/http"

	"github.com/rd-code/kylin/log"
	"github.com/rd-code/kylin/route"
)

type Engine struct {
	router route.Router[func(ctx *Context)]
	logger log.Logger
}

func NewEngine() *Engine {
	return &Engine{
		router: route.NewRouter[func(ctx *Context)](),
		logger: log.Default,
	}
}

func (e *Engine) SetLogger(log log.Logger) {
	e.logger = log
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

func (e *Engine) AddHandler(f ...func(ctx *Context)) {
	e.router.AddHandlers(f...)
}

//group 放到哪里合适
func (e *Engine) Group(path string) route.Router[func(*Context)] {
	return e.router.Group(path)
}

func (e *Engine) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	e.logger.Info("receive request, method:%s, path:%s, ip:%s", request.Method, request.URL.Path, request.RemoteAddr)
	handler := e.router.GetHandler(request.Method, request.URL.Path)
	if handler == nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	ctx := &Context{
		Request:  request,
		Response: response,
		//rs:       handler,
	}
	handler(ctx)
}

func (e *Engine) Listen(addr string) {
	http.ListenAndServe(addr, e)
}
