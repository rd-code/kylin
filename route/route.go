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
	GetHandler(method, path string) Handler
}

type routerImpl struct {
	parent     *routerImpl
	path       string
	methodPath map[string]map[string]Handler
	allPath    map[string]Handler
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
	if _, ok := parent.methodPath[method]; !ok {
		parent.methodPath[method] = make(map[string]Handler)
	}

	data := parent.methodPath[method]
	if _, ok := data[path]; ok {
		panic(fmt.Sprintf("the path:%s is already register", path))
	}
	data [path] = handler
}

func (r *routerImpl) Handle(path string, handler Handler) {
	parent := r.getParent()
	path = addPath(r.path, path)
	if _, ok := parent.allPath[path]; ok {
		panic("the path is already register, path:" + path)
	}
	parent.allPath[path] = handler
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

func (r *routerImpl) GetHandler(method, path string) Handler {
	parent := r.getParent()
	if _, ok := parent.methodPath[method]; ok {
		if v := parent.methodPath[method][path]; v != nil {
			return v
		}
	}
	if v := parent.allPath[path]; v != nil {
		return v
	}
	return nil
}

func NewRouter() Router {
	return &routerImpl{
		methodPath: make(map[string]map[string]Handler),
		allPath:    make(map[string]Handler),
	}
}
