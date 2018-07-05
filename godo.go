package godo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type GoDo struct {
	routes      []*route
	middlewares []func(HandleFunc) HandleFunc
}

type route struct {
	method  string
	pattern string
	handler HandleFunc
}

type HandleFunc func(c *Context)

type Context struct {
	w http.ResponseWriter
	r *http.Request
	*param
}

type param struct {
	form  url.Values
	query url.Values
	body  []byte
	src   string
}

func (p *param) Form(key string) string {
	return p.form.Get(key)
}

func (p *param) Query(key string) string {
	return p.query.Get(key)
}

func (p *param) Src() string {
	return p.src
}

func (c *Context) Bind(v interface{}) error {
	if err := json.Unmarshal(c.param.body, v); err != nil {
		return err
	}
	return nil
}

func New() *GoDo {
	return &GoDo{routes: make([]*route, 0)}
}

func (g *GoDo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimRight(r.URL.Path, "/")
	method := r.Method

	r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	c := &Context{w, r, &param{r.Form, r.URL.Query(), body, ""}}

	index := -1
	for idx, route := range g.routes {
		if route.method != method {
			continue
		}

		if route.pattern == path {
			index = idx
			break
		}

		exist := strings.Contains(route.pattern, ":")

		if exist {

			v1 := strings.Split(route.pattern, "/:")[0]
			arr := strings.Split(path, "/")
			v2 := strings.Join(arr[0:len(arr)-1], "/")

			if v1 == v2 {
				index = idx
				c.param.src = arr[len(arr)-1]
				break
			}
		}
	}
	if index != -1 {
		handler := g.routes[index].handler
		for i := len(g.middlewares) - 1; i >= 0; i-- {
			handler = g.middlewares[i](handler)
		}
		handler(c)
	} else {
		notFoundHandle(c)
	}
}

func notFoundHandle(c *Context) {
	fmt.Fprintf(c.w, "Not Found")
}

func (g *GoDo) Get(pattern string, handler func(c *Context)) {
	g.routes = append(g.routes, &route{"GET", pattern, HandleFunc(handler)})
}

func (g *GoDo) Post(pattern string, handler func(c *Context)) {
	g.routes = append(g.routes, &route{"POST", pattern, HandleFunc(handler)})
}

func (g *GoDo) Put(pattern string, handler func(c *Context)) {
	g.routes = append(g.routes, &route{"PUT", pattern, HandleFunc(handler)})
}

func (g *GoDo) Delete(pattern string, handler func(c *Context)) {
	g.routes = append(g.routes, &route{"DELETE", pattern, HandleFunc(handler)})
}

func (c *Context) JSON(code int, resp interface{}) {
	c.w.Header().Set("content-type", "application/json")
	c.w.WriteHeader(200)
	c.writeJSON(resp)
}

func (c *Context) writeJSON(resp interface{}) {
	j, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	c.w.Write(j)
}

func (c *Context) StatusCode(code int) {
	c.w.WriteHeader(code)
}

func (c *Context) Header() http.Header {
	return c.w.Header()
}

func (g *GoDo) Use(middleware func(HandleFunc) HandleFunc) {
	g.middlewares = append(g.middlewares, middleware)
}

func (g *GoDo) Run(addr ...string) error {
	if len(addr) != 0 {
		return http.ListenAndServe(addr[0], g)
	} else {
		return http.ListenAndServe(":8080", g)
	}
}
