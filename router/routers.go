package router

import (
	"net/http"
	"strings"
	"sync"

	"github.com/kovey/kow/context"
	"github.com/kovey/kow/middleware"
)

type Routers struct {
	middlewares            []context.MiddlewareInterface
	trees                  map[string]*node
	paramsPool             sync.Pool
	maxParams              uint16
	SaveMatchedRoutePath   bool
	RedirectTrailingSlash  bool
	RedirectFixedPath      bool
	HandleMethodNotAllowed bool
	HandleOPTIONS          bool
	GlobalOPTIONS          *Chain
	globalAllowed          string
	NotFound               *Chain
	MethodNotAllowed       *Chain
}

func NewRouters() *Routers {
	return &Routers{trees: make(map[string]*node)}
}

func (r *Routers) getParams() context.Params {
	ps, _ := r.paramsPool.Get().(context.Params)
	return ps
}

func (r *Routers) putParams(ps context.Params) {
	ps.Reset()
	r.paramsPool.Put(ps)
}

func (r *Routers) saveMatchedRoutePath(path string, chain *Chain) *Chain {
	chain.Middlewares = append(chain.Middlewares, middleware.NewSaveMatchedRoute(path))
	return chain
}

func (r *Routers) Handle(method, path string, chain *Chain) {
	varsCount := uint16(0)

	if method == "" {
		panic("method must not be empty")
	}
	if len(path) < 1 || path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}
	if chain == nil {
		panic("chain must not be nil")
	}

	if r.SaveMatchedRoutePath {
		varsCount++
		chain = r.saveMatchedRoutePath(path, chain)
	}

	root := r.trees[method]
	if root == nil {
		root = new(node)
		r.trees[method] = root

		r.globalAllowed = r.allowed("*", "")
	}

	root.addRoute(path, chain)

	if paramsCount := countParams(path); paramsCount+varsCount > r.maxParams {
		r.maxParams = paramsCount + varsCount
	}

	if r.paramsPool.New == nil && r.maxParams > 0 {
		r.paramsPool.New = func() interface{} {
			return make(context.Params, r.maxParams)
		}
	}
}

func (r *Routers) allowed(path, reqMethod string) (allow string) {
	allowed := make([]string, 0, 9)

	if path == "*" {
		if reqMethod == "" {
			for method := range r.trees {
				if method == http.MethodOptions {
					continue
				}
				allowed = append(allowed, method)
			}
		} else {
			return r.globalAllowed
		}
	} else {
		for method := range r.trees {
			if method == reqMethod || method == http.MethodOptions {
				continue
			}

			handle, _, _ := r.trees[method].getValue(path, nil)
			if handle != nil {
				allowed = append(allowed, method)
			}
		}
	}

	if len(allowed) > 0 {
		allowed = append(allowed, http.MethodOptions)

		for i, l := 1, len(allowed); i < l; i++ {
			for j := i; j > 0 && allowed[j] < allowed[j-1]; j-- {
				allowed[j], allowed[j-1] = allowed[j-1], allowed[j]
			}
		}

		return strings.Join(allowed, ", ")
	}

	return allow
}

func (r *Routers) Lookup(method, path string) (*Chain, context.Params, bool) {
	if root := r.trees[method]; root != nil {
		chain, ps, tsr := root.getValue(path, r.getParams)
		if chain == nil {
			r.putParams(ps)
			return nil, nil, tsr
		}
		if ps == nil {
			return chain, nil, tsr
		}
		return chain, ps, tsr
	}
	return nil, nil, false
}

func (r *Routers) HandleHTTP(c *context.Context) {
	path := c.Request.URL.Path
	c.Middleware(r.middlewares...)

	if root := r.trees[c.Request.Method]; root != nil {
		if chain, ps, tsr := root.getValue(path, r.getParams); chain != nil {
			c.Rules = chain.rules
			c.ReqData = chain.param.Clone()
			if ps != nil {
				c.SetParams(ps)
				chain.handle(c)
				r.putParams(ps)
			} else {
				if r.SaveMatchedRoutePath {
					tmpPs := make(context.Params)
					c.SetParams(tmpPs)
				}
				chain.handle(c)
			}
			return
		} else if c.Request.Method != http.MethodConnect && path != "/" {
			code := http.StatusMovedPermanently
			if c.Request.Method != http.MethodGet {
				code = http.StatusPermanentRedirect
			}

			if tsr && r.RedirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					c.Request.URL.Path = path[:len(path)-1]
				} else {
					c.Request.URL.Path = path + "/"
				}
				http.Redirect(c.Writer(), c.Request, c.Request.URL.String(), code)
				return
			}

			if r.RedirectFixedPath {
				fixedPath, found := root.findCaseInsensitivePath(
					CleanPath(path),
					r.RedirectTrailingSlash,
				)
				if found {
					c.Request.URL.Path = fixedPath
					http.Redirect(c.Writer(), c.Request, c.Request.URL.String(), code)
					return
				}
			}
		}
	}

	if c.Request.Method == http.MethodOptions && r.HandleOPTIONS {
		if allow := r.allowed(path, http.MethodOptions); allow != "" {
			c.Header("Allow", allow)
			if r.GlobalOPTIONS != nil {
				r.GlobalOPTIONS.handle(c)
			}
			return
		}
	} else if r.HandleMethodNotAllowed {
		if allow := r.allowed(path, c.Request.Method); allow != "" {
			c.Header("Allow", allow)
			if r.MethodNotAllowed != nil {
				r.MethodNotAllowed.handle(c)
			} else {
				http.Error(c.Writer(),
					http.StatusText(http.StatusMethodNotAllowed),
					http.StatusMethodNotAllowed,
				)
			}
			return
		}
	}

	if r.NotFound != nil {
		r.NotFound.handle(c)
	} else {
		http.NotFound(c.Writer(), c.Request)
	}
}

func (r *Routers) Middlerware(middlewares ...context.MiddlewareInterface) {
	r.middlewares = append(r.middlewares, middlewares...)
}

func (r *Routers) Add(router RouterInterface) {
	r.Handle(router.Method(), router.Path(), router.Chain())
}

func (r *Routers) ServeFiles(path string, root http.FileSystem, chain *Chain) {
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path must end with /*filepath in path '" + path + "'")
	}

	chain.fileServer = http.FileServer(root)
	r.Handle(http.MethodGet, path, chain)
}

func (r *Routers) Group(path string) *Group {
	return NewGroup(path, r)
}
