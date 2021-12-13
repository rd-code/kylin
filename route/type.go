package route

type paramElement struct {
	start int
	end   int
	key   string
}

type customPath struct {
	path     string
	elements []*paramElement
	common   bool
}

func (c *customPath) match(path string) (res map[string]string, ok bool) {
	if c.common {
		ok = c.path == path
		return
	}
	var pCount, tCount int
	res = make(map[string]string)
	for _, element := range c.elements {
		size := element.start - tCount - 2

		if path[pCount:pCount+size] != c.path[tCount:tCount+size] {
			return
		}
		pCount += size + 1
		tCount = element.end
		var index int
		for index = pCount; index < len(path); index++ {
			if path[index] == Slash {
				res[element.key] = path[pCount:index]
				pCount = index
				break
			}
		}
		if index == len(path) {
			res[element.key] = path[pCount:index]
			pCount = index
		}
	}
	if pCount== len(path) && tCount== len(c.path){
		ok = true
		return
	}
	if pCount == len(path) || tCount == len(c.path){
		return
	}

	if path[pCount:] == c.path[tCount:]{
		ok = true
	}
	return
}
