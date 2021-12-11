package route

import (
	"fmt"
	"testing"
)

func TestFormatPath(t *testing.T) {
	data := map[string]string{
		"abc":   "/abc",
		"/abc":  "/abc",
		"abc/":  "/abc",
		"/abc/": "/abc",
	}
	for k, v := range data {
		value := formatPath(k)
		if v != value {
			t.Fatalf("format path:%s failed, expected:%s, acutal:%s",
				k, v, value)
		}
	}
}

func TestAddPath(t *testing.T) {
	type data struct {
		path1 string
		path2 string
		res   string
	}
	var items = []data{
		{
			path1: "/abc",
			path2: "/bbc",
			res: "/abc/bbc",
		},
		{
			path1: "/abc/",
			path2: "/bbc/",
			res: "/abc/bbc",
		},
		{
			path1: "/",
			path2: "/abc",
			res: "/abc",
		},
		{
			path1: "/abc",
			path2: "/",
			res: "/abc",
		},
		{
			path1: "",
			path2: "",
			res: "/",
		},
	}
	for _, item := range items {
		value := addPath(item.path1, item.path2)
		if value != item.res {
			t.Fatalf("add path failed, input paht1:%s path2:%s, expected:%s, actual:%s",
				item.path1, item.path2, item.res, value)
		}
	}
}

func TestGetCustomPath(t *testing.T) {
	data := make(map[string]*customPath)
	_ = data
}

func TestParse(t *testing.T) {
	c := getCustomPath("/abc/:bce/ccc/:cde/a")
	v, ok := c.match("/abc/张三/ccc/aaaa")
	fmt.Println(ok)
	for k, v := range v {
		fmt.Println(k, v)
	}
}

func TestMatch(t *testing.T) {

}
