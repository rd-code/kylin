package route

import (
	"path"
	"strings"
)

//格式化path, 最终输出为/api/v2/v1
func formatPath(url string) string {
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}
	if len(url) != 1 && url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}
	path.Clean(url)
	return url
}

//将两个path加起来,path1在前面，path2在后面
// eg path1:/api/v1 path2:/user/test res:/api/v1/user/test
func addPath(path1, path2 string) string {
	path1 = formatPath(path1)
	path2 = formatPath(path2)
	return path.Join(path1, path2)
}

//检查path是否合法
func checkPath(path string) bool {
	if len(path) == 0 {
		return false
	}
	if len(path) == 1 {
		return path == "/"
	}
	return strings.HasPrefix(path, "/") && strings.HasSuffix(path, "/")
}

func isCommon(path string) bool {
	return strings.Contains(path, "/:")
}
