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
	GetHandler(method, path string) *RouterServer[T]
}

type routerHandler[T any] struct {
	handler    T
	customPath *customPath
}

type routerMux[T any] struct {
	common map[string]*routerHandler[T]
	//	handlers []*routerHandler[T]
}

func (r *routerMux[T]) addCommon(path string, t T) {
	if _, ok := r.common[path]; ok {
		panic(fmt.Sprintf("the path:%s is already register", path))
	}
	cp := getCustomPath(path)
	/*	if !cp.common {
		panic(fmt.Sprintf("the path:%s is not common", path))
	}*/
	r.common[path] = &routerHandler[T]{
		handler:    t,
		customPath: cp,
	}
}

func (r *routerMux[T]) add(path string, t T) {
	/*for _, h := range r.handlers {
		if h.customPath.path == path {
			panic(fmt.Sprintf("the path:%s is already register", path))
		}
	}*/
	cp := getCustomPath(path)
	if cp.common {
		panic(fmt.Sprintf("the path:%s is common", path))
	}
	/*r.handlers = append(r.handlers, &routerHandler[T]{
		handler:    t,
		customPath: cp,
	})*/
}

func (r *routerMux[T]) register(path string, t T) {
	if isCommon(path) {
		r.addCommon(path, t)
		return
	}
	r.add(path, t)
}

func (r *routerMux[T]) getServer(path string) (rs *RouterServer[T]) {
	if v, ok := r.common[path]; ok {
		rs = &RouterServer[T]{
			handler: v,
		}
		return
	}
	/*for _, h := range r.handlers {
		params, ok := h.customPath.match(path)
		if ok {
			rs = &RouterServer{
				params:  params,
				handler: h,
			}
			return
		}
	}*/
	return
}

type RouterServer[T any] struct {
	params  map[string]string
	handler *routerHandler[T]
}

func (r *RouterServer[T]) Handler() T {
	return r.handler.handler
}

func (r *RouterServer) Param(key string) string {
	return r.params[key]
}

type routerImpl[T any] struct {
	parent    *routerImpl[T]
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
	parent := r.getParent()
	_ = parent

	path = addPath(r.path, path)

	return &routerImpl[T]{
		parent: parent,
		path:   path,
	}

}

func (r *routerImpl[T]) register(method, path string, t T) {
	parent := r.getParent()
	path = addPath(r.path, path)
	if _, ok := parent.routerMux[method]; !ok {
		parent.routerMux[method] = &routerMux[T]{}
	}

	rm := parent.routerMux[method]
	rm.register(path, t)

}

const (
	MethodHandle = "handle"
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

func (r *routerImpl[T]) getHandler(method, path string) (rs *RouterServer[T]) {
	parent := r.getParent()
	if v, ok := parent.routerMux[method]; ok {
		rs = v.getServer(path)
		if rs != nil {
			return
		}
	}
	return
}

func (r *routerImpl[T]) GetHandler(method, path string) (rs *RouterServer[T]) {
	if rs = r.getHandler(method, path); rs != nil {
		return
	}
	rs = r.getHandler(MethodHandle, path)
	return
}

func NewRouter[T any]() Router[T] {
	return &routerImpl[T]{
		routerMux: make(map[string]*routerMux[T]),
	}
}
