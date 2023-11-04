package router

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/kovey/kow/context"
)

func min(a, b int) int {
	if a > b {
		return b
	}

	return a
}

func longestCommonPrefix(a, b string) int {
	i := 0
	max := min(len(a), len(b))
	for i < max && a[i] == b[i] {
		i++
	}

	return i
}

func findWildcard(path string) (string, int, bool) {
	for start, c := range []byte(path) {
		if c != ':' && c != '*' {
			continue
		}

		valid := true

		for end, c := range []byte(path[start+1:]) {
			switch c {
			case '/':
				return path[start : start+1+end], start, valid
			case ':', '*':
				valid = false
			}
		}

		return path[start:], start, valid
	}

	return "", -1, false
}

func countParams(path string) uint16 {
	var n uint16
	for _, p := range []byte(path) {
		switch p {
		case ':', '*':
			n++
		}
	}

	return n
}

type nodeType uint8

const (
	static nodeType = iota
	root
	param
	catchAll
)

type node struct {
	path      string
	indices   string
	wildChild bool
	nType     nodeType
	priority  uint32
	children  []*node
	chain     *Chain
}

func (n *node) incrementChildPrio(pos int) int {
	cs := n.children
	cs[pos].priority++
	prio := cs[pos].priority

	newPos := pos
	for ; newPos > 0 && cs[newPos-1].priority < prio; newPos-- {
		cs[newPos-1], cs[newPos] = cs[newPos], cs[newPos-1]
	}

	if newPos != pos {
		n.indices = n.indices[:newPos] +
			n.indices[pos:pos+1] +
			n.indices[newPos:pos] + n.indices[pos+1:]
	}

	return newPos
}
func (n *node) addRoute(path string, chain *Chain) {
	fullPath := path
	n.priority++

	if n.path == "" && n.indices == "" {
		n.insertChild(path, fullPath, chain)
		n.nType = root
		return
	}

walk:
	for {
		i := longestCommonPrefix(path, n.path)

		if i < len(n.path) {
			child := node{
				path:      n.path[i:],
				wildChild: n.wildChild,
				nType:     static,
				indices:   n.indices,
				children:  n.children,
				chain:     n.chain,
				priority:  n.priority - 1,
			}

			n.children = []*node{&child}
			n.indices = string([]byte{n.path[i]})
			n.path = path[:i]
			n.chain = nil
			n.wildChild = false
		}

		if i < len(path) {
			path = path[i:]

			if n.wildChild {
				n = n.children[0]
				n.priority++

				if len(path) >= len(n.path) && n.path == path[:len(n.path)] &&
					n.nType != catchAll &&
					(len(n.path) >= len(path) || path[len(n.path)] == '/') {
					continue walk
				} else {
					pathSeg := path
					if n.nType != catchAll {
						pathSeg = strings.SplitN(pathSeg, "/", 2)[0]
					}
					prefix := fullPath[:strings.Index(fullPath, pathSeg)] + n.path
					panic("'" + pathSeg +
						"' in new path '" + fullPath +
						"' conflicts with existing wildcard '" + n.path +
						"' in existing prefix '" + prefix +
						"'")
				}
			}

			idxc := path[0]

			if n.nType == param && idxc == '/' && len(n.children) == 1 {
				n = n.children[0]
				n.priority++
				continue walk
			}

			for i, c := range []byte(n.indices) {
				if c == idxc {
					i = n.incrementChildPrio(i)
					n = n.children[i]
					continue walk
				}
			}

			if idxc != ':' && idxc != '*' {
				n.indices += string([]byte{idxc})
				child := &node{}
				n.children = append(n.children, child)
				n.incrementChildPrio(len(n.indices) - 1)
				n = child
			}
			n.insertChild(path, fullPath, chain)
			return
		}

		if n.chain != nil {
			panic("a handle is already registered for path '" + fullPath + "'")
		}
		n.chain = chain
		return
	}
}

func (n *node) insertChild(path, fullPath string, chain *Chain) {
	for {
		wildcard, i, valid := findWildcard(path)
		if i < 0 {
			break
		}

		if !valid {
			panic("only one wildcard per path segment is allowed, has: '" +
				wildcard + "' in path '" + fullPath + "'")
		}

		if len(wildcard) < 2 {
			panic("wildcards must be named with a non-empty name in path '" + fullPath + "'")
		}

		if len(n.children) > 0 {
			panic("wildcard segment '" + wildcard +
				"' conflicts with existing children in path '" + fullPath + "'")
		}

		if wildcard[0] == ':' {
			if i > 0 {
				n.path = path[:i]
				path = path[i:]
			}

			n.wildChild = true
			child := &node{
				nType: param,
				path:  wildcard,
			}
			n.children = []*node{child}
			n = child
			n.priority++

			if len(wildcard) < len(path) {
				path = path[len(wildcard):]
				child := &node{
					priority: 1,
				}
				n.children = []*node{child}
				n = child
				continue
			}

			n.chain = chain
			return
		}

		if i+len(wildcard) != len(path) {
			panic("catch-all routes are only allowed at the end of the path in path '" + fullPath + "'")
		}

		if len(n.path) > 0 && n.path[len(n.path)-1] == '/' {
			panic("catch-all conflicts with existing handle for the path segment root in path '" + fullPath + "'")
		}

		i--
		if path[i] != '/' {
			panic("no / before catch-all in path '" + fullPath + "'")
		}

		n.path = path[:i]

		child := &node{
			wildChild: true,
			nType:     catchAll,
		}
		n.children = []*node{child}
		n.indices = string('/')
		n = child
		n.priority++

		child = &node{
			path:     path[i:],
			nType:    catchAll,
			chain:    chain,
			priority: 1,
		}
		n.children = []*node{child}

		return
	}

	n.path = path
	n.chain = chain
}

func (n *node) getValue(path string, params func() context.Params) (chain *Chain, ps context.Params, tsr bool) {
walk:
	for {
		prefix := n.path
		if len(path) > len(prefix) {
			if path[:len(prefix)] == prefix {
				path = path[len(prefix):]

				if !n.wildChild {
					idxc := path[0]
					for i, c := range []byte(n.indices) {
						if c == idxc {
							n = n.children[i]
							continue walk
						}
					}

					tsr = (path == "/" && n.chain != nil)
					return
				}

				n = n.children[0]
				switch n.nType {
				case param:
					end := 0
					for end < len(path) && path[end] != '/' {
						end++
					}

					if params != nil {
						if ps == nil {
							ps = params()
						}
						ps[n.path[1:]] = path[:end]
					}

					if end < len(path) {
						if len(n.children) > 0 {
							path = path[end:]
							n = n.children[0]
							continue walk
						}

						tsr = (len(path) == end+1)
						return
					}

					if chain = n.chain; chain != nil {
						return
					} else if len(n.children) == 1 {
						n = n.children[0]
						tsr = (n.path == "/" && n.chain != nil) || (n.path == "" && n.indices == "/")
					}

					return

				case catchAll:
					if params != nil {
						if ps == nil {
							ps = params()
						}
						ps[n.path[2:]] = path
					}

					chain = n.chain
					return

				default:
					panic("invalid node type")
				}
			}
		} else if path == prefix {
			if chain = n.chain; chain != nil {
				return
			}

			if path == "/" && n.wildChild && n.nType != root {
				tsr = true
				return
			}

			if path == "/" && n.nType == static {
				tsr = true
				return
			}

			for i, c := range []byte(n.indices) {
				if c == '/' {
					n = n.children[i]
					tsr = (len(n.path) == 1 && n.chain != nil) ||
						(n.nType == catchAll && n.children[0].chain != nil)
					return
				}
			}
			return
		}

		tsr = (path == "/") ||
			(len(prefix) == len(path)+1 && prefix[len(path)] == '/' &&
				path == prefix[:len(prefix)-1] && n.chain != nil)
		return
	}
}

func (n *node) findCaseInsensitivePath(path string, fixTrailingSlash bool) (fixedPath string, found bool) {
	const stackBufSize = 128

	buf := make([]byte, 0, stackBufSize)
	if l := len(path) + 1; l > stackBufSize {
		buf = make([]byte, 0, l)
	}

	ciPath := n.findCaseInsensitivePathRec(
		path,
		buf,
		[4]byte{},
		fixTrailingSlash,
	)

	return string(ciPath), ciPath != nil
}

func shiftNRuneBytes(rb [4]byte, n int) [4]byte {
	switch n {
	case 0:
		return rb
	case 1:
		return [4]byte{rb[1], rb[2], rb[3], 0}
	case 2:
		return [4]byte{rb[2], rb[3]}
	case 3:
		return [4]byte{rb[3]}
	default:
		return [4]byte{}
	}
}

func (n *node) findCaseInsensitivePathRec(path string, ciPath []byte, rb [4]byte, fixTrailingSlash bool) []byte {
	npLen := len(n.path)

walk:
	for len(path) >= npLen && (npLen == 0 || strings.EqualFold(path[1:npLen], n.path[1:])) {
		oldPath := path
		path = path[npLen:]
		ciPath = append(ciPath, n.path...)

		if len(path) > 0 {
			if !n.wildChild {
				rb = shiftNRuneBytes(rb, npLen)

				if rb[0] != 0 {
					idxc := rb[0]
					for i, c := range []byte(n.indices) {
						if c == idxc {
							n = n.children[i]
							npLen = len(n.path)
							continue walk
						}
					}
				} else {
					var rv rune

					var off int
					for max := min(npLen, 3); off < max; off++ {
						if i := npLen - off; utf8.RuneStart(oldPath[i]) {
							rv, _ = utf8.DecodeRuneInString(oldPath[i:])
							break
						}
					}

					lo := unicode.ToLower(rv)
					utf8.EncodeRune(rb[:], lo)

					rb = shiftNRuneBytes(rb, off)

					idxc := rb[0]
					for i, c := range []byte(n.indices) {
						if c == idxc {
							if out := n.children[i].findCaseInsensitivePathRec(
								path, ciPath, rb, fixTrailingSlash,
							); out != nil {
								return out
							}
							break
						}
					}

					if up := unicode.ToUpper(rv); up != lo {
						utf8.EncodeRune(rb[:], up)
						rb = shiftNRuneBytes(rb, off)

						idxc := rb[0]
						for i, c := range []byte(n.indices) {
							if c == idxc {
								n = n.children[i]
								npLen = len(n.path)
								continue walk
							}
						}
					}
				}

				if fixTrailingSlash && path == "/" && n.chain != nil {
					return ciPath
				}
				return nil
			}

			n = n.children[0]
			switch n.nType {
			case param:
				end := 0
				for end < len(path) && path[end] != '/' {
					end++
				}

				ciPath = append(ciPath, path[:end]...)

				if end < len(path) {
					if len(n.children) > 0 {
						n = n.children[0]
						npLen = len(n.path)
						path = path[end:]
						continue
					}

					if fixTrailingSlash && len(path) == end+1 {
						return ciPath
					}
					return nil
				}

				if n.chain != nil {
					return ciPath
				} else if fixTrailingSlash && len(n.children) == 1 {
					n = n.children[0]
					if n.path == "/" && n.chain != nil {
						return append(ciPath, '/')
					}
				}
				return nil

			case catchAll:
				return append(ciPath, path...)

			default:
				panic("invalid node type")
			}
		} else {
			if n.chain != nil {
				return ciPath
			}

			if fixTrailingSlash {
				for i, c := range []byte(n.indices) {
					if c == '/' {
						n = n.children[i]
						if (len(n.path) == 1 && n.chain != nil) ||
							(n.nType == catchAll && n.children[0].chain != nil) {
							return append(ciPath, '/')
						}
						return nil
					}
				}
			}
			return nil
		}
	}

	if fixTrailingSlash {
		if path == "/" {
			return ciPath
		}
		if len(path)+1 == npLen && n.path[len(path)] == '/' &&
			strings.EqualFold(path[1:], n.path[1:len(path)]) && n.chain != nil {
			return append(ciPath, n.path...)
		}
	}
	return nil
}
