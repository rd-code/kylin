package route

import (
	"fmt"
	"net/http"
)

type Handler interface {
}

type Router interface {
	Get(path string, handler Handler)
	Post(path string, handler Handler)
	Delete(path string, handler Handler)
	Put(path string, handler Handler)
	Handle(path string, handler Handler)
	Group(path string) Router
	GetHandler(method, path string) *RouterServer
}

type routerHandler struct {
	handler    Handler
	customPath *customPath
}

type routerMux struct {
	common   map[string]*routerHandler
	handlers []*routerHandler
}

func (r *routerMux) addCommon(path string, handler Handler) {
	if _, ok := r.common[path]; ok {
		panic(fmt.Sprintf("the path:%s is already register", path))
	}
	cp := getCustomPath(path)
	if !cp.common {
		panic(fmt.Sprintf("the path:%s is not common", path))
	}
	r.common[path] = &routerHandler{
		handler:    handler,
		customPath: cp,
	}
}

func (r *routerMux) add(path string, handler Handler) {
	for _, h := range r.handlers {
		if h.customPath.path == path {
			panic(fmt.Sprintf("the path:%s is already register", path))
		}
	}
	cp := getCustomPath(path)
	if !cp.common {
		panic(fmt.Sprintf("the path:%s is common", path))
	}
	r.handlers = append(r.handlers, &routerHandler{
		handler:    handler,
		customPath: cp,
	})
}

func (r *routerMux) register(path string, handler Handler) {
	if isCommon(path) {
		r.addCommon(path, handler)
		return
	}
	r.add(path, handler)
}

func (r *routerMux) getServer(path string) (rs *RouterServer) {
	if v, ok := r.common[path]; ok {
		rs = &RouterServer{
			handler: v,
		}
		return
	}
	for _, h := range r.handlers {
		params, ok := h.customPath.match(path)
		if ok {
			rs = &RouterServer{
				params:  params,
				handler: h,
			}
			return
		}
	}
	return
}

type RouterServer struct {
	params  map[string]string
	handler *routerHandler
}

func (r *RouterServer) Handler() Handler {
	return &r.handler.handler
}

func (r *RouterServer) Param(key string) string {
	return r.params[key]
}

type routerImpl struct {
	parent    *routerImpl
	path      string
	routerMux map[string]*routerMux
}

func (r *routerImpl) getParent() *routerImpl {
	parent := r.parent
	if parent == nil {
		parent = r
	}
	return parent
}

func (r *routerImpl) Group(path string) Router {
	parent := r.getParent()

	path = addPath(r.path, path)

	return &routerImpl{
		parent: parent,
		path:   path,
	}
}

func (r *routerImpl) register(method, path string, handler Handler) {
	parent := r.getParent()
	path = addPath(r.path, path)
	if _, ok := parent.routerMux[method]; !ok {
		parent.routerMux[method] = &routerMux{}
	}

	rm := parent.routerMux[method]
	rm.register(path, handler)

}

const (
	MethodHandle = "handle"
)

func (r *routerImpl) Handle(path string, handler Handler) {
	r.register(MethodHandle, path, handler)

}

func (r *routerImpl) Get(path string, handler Handler) {
	r.register(http.MethodGet, path, handler)
}

func (r *routerImpl) Post(path string, handler Handler) {
	r.register(http.MethodPost, path, handler)
}

func (r *routerImpl) Delete(path string, handler Handler) {
	r.register(http.MethodDelete, path, handler)
}

func (r *routerImpl) Put(path string, handler Handler) {
	r.register(http.MethodPut, path, handler)
}

func (r *routerImpl) getHandler(method, path string) (rs *RouterServer) {
	parent := r.getParent()
	if v, ok := parent.routerMux[method]; ok {
		rs = v.getServer(path)
		if rs != nil {
			return
		}
	}
	return
}

func (r *routerImpl) GetHandler(method, path string) (rs *RouterServer) {
	if rs = r.getHandler(method, path); rs != nil {
		return
	}
	rs = r.getHandler(MethodHandle, path)
	return
}

func NewRouter() Router {
	return &routerImpl{
		routerMux: make(map[string]*routerMux),
	}
}
