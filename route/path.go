package route

import (
	"fmt"
	"strings"
)

//格式化path, 最终输出为/api/v2/v1
func formatPath(path string) string {
	if len(path) == 0 {
		panic("the path is empty")
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if len(path) != 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return path
}

//将两个path加起来,path1在前面，path2在后面
// eg path1:/api/v1 path2:/user/test res:/api/v1/user/test
func addPath(path1, path2 string) string {
	path1 = formatPath(path1)
	path2 = formatPath(path2)
	if len(path1) == 1 {
		return path2
	}
	if len(path2) == 1 {
		return path1
	}
	return path1 + path2
}

func isCommon(path string) bool {
	return !strings.Contains(path, "/:")
}

const (
	Slash byte = 47
	Colon byte = 58
)

func getParams(path *string, start, end int) *paramElement {
	if start >= end {
		panic(fmt.Errorf("get params: the path:%s is invalid", *path))
	}
	return &paramElement{
		start: start,
		end:   end,
		key:   (*path)[start:end],
	}
}

// 解析path
func parse(path string) []*paramElement {
	if isCommon(path) {
		return nil
	}
	size := len(path)

	flag := false
	var keyStart int
	var params []*paramElement
	for i := 0; i < size; i++ {
		v := path[i]
		if v == Slash {
			if i+1 < size && path[i+1] == Colon {
				if flag {
					panic("the path is invalid:" + path)
				}
				flag = true
				i += 1
				keyStart = i
			} else {
				if flag {
					params = append(params, getParams(&path, keyStart+1, i))
					flag = false
				}
			}
		}
	}
	if flag {
		params = append(params, getParams(&path, keyStart+1, size))
	}
	return params
}

func getCustomPath(path string) *customPath {
	items := parse(path)
	return &customPath{
		path:     path,
		elements: items,
		common:   items == nil,
	}
}
