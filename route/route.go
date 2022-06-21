package route

import (
	"fmt"
	"net/http"
)

type Router[T any] interface {
	Get(path string, t T)
	Post(path string, t T)
	Delete(path string, t T)
	Put(path string, t T)
	Handle(path string, t T)
	Group(path string) Router[T]
	GetHandler(method, path string) T
	AddHandlers(t ...T)
}

//存储某个method下的所有请求和注册的handler
type routerMux[T any] struct {
	common map[string]T
}

func (r *routerMux[T]) addCommon(path string, t T) {
	if _, ok := r.common[path]; ok {
		panic(fmt.Sprintf("the path:%s is already register", path))
	}
	r.common[path] = t
}

func (r *routerMux[T]) register(path string, t T) {
	if isCommon(path) {
		r.addCommon(path, t)
		return
	}
}

type routerImpl[T any] struct {
	children  []*routerImpl[T]
	parent    *routerImpl[T]
	handlers  []T
	path      string
	routerMux map[string]*routerMux[T]
}

func (r *routerImpl[T]) getParent() *routerImpl[T] {
	parent := r.parent
	if parent == nil {
		parent = r
	}
	return parent
}

func (r *routerImpl[T]) Group(path string) Router[T] {
	if len(path) == 0 || path == "/" {
		panic("the path:" + path + " is invalid")
	}
	parent := r.getParent()

	path = addPath(r.path, path)

	return &routerImpl[T]{
		parent: parent,
		path:   path,
	}

}

func (r *routerImpl[T]) AddHandlers(t ...T) {
	r.handlers = append(r.handlers, t...)
}

func (r *routerImpl[T]) register(method, path string, t T) {
	parent := r.getParent()
	path = addPath(r.path, path)
	if _, ok := parent.routerMux[method]; !ok {
		parent.routerMux[method] = &routerMux[T]{
			common: make(map[string]T),
		}
	}

	rm := parent.routerMux[method]
	rm.register(path, t)

}

const (
	//未指定method的请求， 可以匹配任何method
	MethodHandle = "HANDLE"
)

func (r *routerImpl[T]) Handle(path string, handler T) {
	r.register(MethodHandle, path, handler)

}

func (r *routerImpl[T]) Get(path string, handler T) {
	r.register(http.MethodGet, path, handler)
}

func (r *routerImpl[T]) Post(path string, handler T) {
	r.register(http.MethodPost, path, handler)
}

func (r *routerImpl[T]) Delete(path string, handler T) {
	r.register(http.MethodDelete, path, handler)
}

func (r *routerImpl[T]) Put(path string, handler T) {
	r.register(http.MethodPut, path, handler)
}

func (r *routerImpl[T]) getHandler(method, path string) (t T) {
	parent := r.getParent()

	if v, ok := parent.routerMux[method]; ok {
		t = v.common[path]
	}
	return
}

func (r *routerImpl[T]) GetHandler(method, path string) (t T) {
	t = r.getHandler(method, path)
	return
}

func NewRouter[T any]() Router[T] {
	return &routerImpl[T]{
		routerMux: make(map[string]*routerMux[T]),
	}
}
