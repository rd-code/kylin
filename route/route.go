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
	GetHandlers(method, path string) []T
	AddHandlers(t ...T)
}

//存储某个method下的所有请求和注册的handler
type routerMux[T any] struct {
	common map[string]map[string][]T
}

func newRouterMux[T any]() *routerMux[T] {
	return &routerMux[T]{
		common: make(map[string]map[string][]T),
	}
}

func (r *routerMux[T]) addCommon(method, path string, t ...T) {
	if _, ok := r.common[method][path]; ok {
		panic(fmt.Sprintf("the path:%s is already register", path))
	}
	r.common[method][path] = t
}

func (r *routerMux[T]) register(method, path string, t ...T) {
	if _, ok := r.common[method]; !ok {
		r.common[method] = make(map[string][]T)
	}

	if isCommon(path) {
		r.addCommon(method, path, t...)
		return
	}
	panic(fmt.Sprintf("the path:%s is invalid", path))
}

type routerImpl[T any] struct {
	children  []*routerImpl[T]
	handlers  []T
	path      string
	routerMux *routerMux[T]
}

//追究handler
func (r *routerImpl[T]) appendHandler(t ...T) (res []T) {
	res = make([]T, 0, len(r.handlers)+len(t))
	res = append(res, r.handlers...)
	res = append(res, t...)
	return
}

func (r *routerImpl[T]) gen(path string) *routerImpl[T] {
	if len(path) == 0 || path == "/" {
		panic("the path:" + path + " is invalid")
	}
	path = addPath(r.path, path)
	res := &routerImpl[T]{
		path:      path,
		handlers:  r.handlers,
		routerMux: r.routerMux,
	}
	r.children = append(r.children, res)
	return res
}

func (r *routerImpl[T]) Group(path string) Router[T] {
	if len(path) == 0 || path == "/" {
		panic("the path:" + path + " is invalid")
	}
	return r.gen(path)
}

func (r *routerImpl[T]) AddHandlers(t ...T) {
	r.handlers = r.appendHandler(t...)
}

func (r *routerImpl[T]) register(method, path string, t T) {
	path = addPath(r.path, path)
	r.routerMux.register(method, path, r.appendHandler(t))
}

func (r *routerImpl[T]) Register(method, path string, t T) {
	r.register(method, path, t)
}

func (r *routerImpl[T]) Handle(path string, handler T) {
	r.register(http.MethodGet, path, handler)
	r.register(http.MethodHead, path, handler)
	r.register(http.MethodPost, path, handler)
	r.register(http.MethodPut, path, handler)
	r.register(http.MethodGet, path, handler)
	r.register(http.MethodPatch, path, handler)
	r.register(http.MethodDelete, path, handler)
	r.register(http.MethodConnect, path, handler)
	r.register(http.MethodOptions, path, handler)
	r.register(http.MethodTrace, path, handler)
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

func (r *routerImpl[T]) getHandlers(method, path string) (t []T) {
	if _, ok := r.routerMux.common[method]; ok {
		return
	}
	t = r.routerMux.common[method][path]
	return
}

func (r *routerImpl[T]) GetHandlers(method, path string) (t []T) {
	t = r.getHandlers(method, path)
	return
}

func NewRouter[T any]() Router[T] {
	return &routerImpl[T]{
		routerMux: newRouterMux[T](),
	}
}
