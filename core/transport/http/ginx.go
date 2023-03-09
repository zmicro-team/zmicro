package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Context struct {
	ctx *gin.Context
}

func (c *Context) BindXxx() {

}

func (c *Context) String(code int, format string, values ...any) {
	c.ctx.String(code, format, values...)
}

type HandlerFunc func(*Context)

func (s *Server) GroupEx(path string, handlers ...HandlerFunc) *Server {
	if len(handlers) == 0 {
		return s
	}
	gHandlers := make([]gin.HandlerFunc, 0, len(handlers))
	for _, h := range handlers {
		handler := func(c *gin.Context) {
			h(&Context{ctx: c})
		}
		gHandlers = append(gHandlers, handler)
	}

	s.Group(path, gHandlers...)

	return s
}

func (s *Server) UseEx(middlewares ...HandlerFunc) {
	if len(middlewares) == 0 {
		return
	}
	handlers := make([]gin.HandlerFunc, 0, len(middlewares))
	for _, h := range middlewares {
		handler := func(c *gin.Context) {
			h(&Context{ctx: c})
		}
		handlers = append(handlers, handler)
	}

	s.Use(handlers...)
}

func (s *Server) HandleEx(method, path string, handlers ...HandlerFunc) *Server {
	if len(handlers) == 0 {
		return s
	}
	gHandlers := make([]gin.HandlerFunc, 0, len(handlers))
	for _, h := range handlers {
		handler := func(c *gin.Context) {
			h(&Context{ctx: c})
		}
		gHandlers = append(gHandlers, handler)
	}

	s.Handle(method, path, gHandlers...)
	return s
}

func (s *Server) PostEx(path string, handlers ...HandlerFunc) *Server {
	return s.HandleEx(http.MethodPost, path, handlers...)
}

func (s *Server) GetEx(path string, handlers ...HandlerFunc) *Server {
	return s.HandleEx(http.MethodGet, path, handlers...)
}

func (s *Server) DeleteEx(path string, handlers ...HandlerFunc) *Server {
	return s.HandleEx(http.MethodDelete, path, handlers...)
}

func (s *Server) PatchEx(path string, handlers ...HandlerFunc) *Server {
	return s.HandleEx(http.MethodPatch, path, handlers...)
}

func (s *Server) PutEx(path string, handlers ...HandlerFunc) *Server {
	return s.HandleEx(http.MethodPut, path, handlers...)
}

func (s *Server) OptionsEx(path string, handlers ...HandlerFunc) *Server {
	return s.HandleEx(http.MethodOptions, path, handlers...)
}

func (s *Server) HeadEx(path string, handlers ...HandlerFunc) *Server {
	return s.HandleEx(http.MethodHead, path, handlers...)
}

func (s *Server) AnyEx(path string, handlers ...HandlerFunc) *Server {
	anyMethods := []string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodConnect,
		http.MethodTrace,
	}
	for _, method := range anyMethods {
		s.HandleEx(method, path, handlers...)
	}

	return s
}

func (s *Server) MatchEx(methods []string, path string, handlers ...HandlerFunc) *Server {
	for _, method := range methods {
		s.HandleEx(method, path, handlers...)
	}

	return s
}
