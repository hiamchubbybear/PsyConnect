



package gin

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/internal/bytesconv"
	"github.com/gin-gonic/gin/render"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const defaultMultipartMemory = 32 << 20 

var (
	default404Body = []byte("404 page not found")
	default405Body = []byte("405 method not allowed")
)

var defaultPlatform string

var defaultTrustedCIDRs = []*net.IPNet{
	{ 
		IP:   net.IP{0x0, 0x0, 0x0, 0x0},
		Mask: net.IPMask{0x0, 0x0, 0x0, 0x0},
	},
	{ 
		IP:   net.IP{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		Mask: net.IPMask{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	},
}

var regSafePrefix = regexp.MustCompile("[^a-zA-Z0-9/-]+")
var regRemoveRepeatedChar = regexp.MustCompile("/{2,}")


type HandlerFunc func(*Context)


type OptionFunc func(*Engine)


type HandlersChain []HandlerFunc


func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}


type RouteInfo struct {
	Method      string
	Path        string
	Handler     string
	HandlerFunc HandlerFunc
}


type RoutesInfo []RouteInfo


const (
	
	
	PlatformGoogleAppEngine = "X-Appengine-Remote-Addr"
	
	
	PlatformCloudflare = "CF-Connecting-IP"
	
	PlatformFlyIO = "Fly-Client-IP"
)



type Engine struct {
	RouterGroup

	
	
	
	
	
	RedirectTrailingSlash bool

	
	
	
	
	
	
	
	
	
	RedirectFixedPath bool

	
	
	
	
	
	
	HandleMethodNotAllowed bool

	
	
	
	
	ForwardedByClientIP bool

	
	
	
	
	AppEngine bool

	
	UseRawPath bool

	
	
	
	UnescapePathValues bool

	
	
	RemoveExtraSlash bool

	
	
	
	
	RemoteIPHeaders []string

	
	
	TrustedPlatform string

	
	
	MaxMultipartMemory int64

	
	UseH2C bool

	
	ContextWithFallback bool

	delims           render.Delims
	secureJSONPrefix string
	HTMLRender       render.HTMLRender
	FuncMap          template.FuncMap
	allNoRoute       HandlersChain
	allNoMethod      HandlersChain
	noRoute          HandlersChain
	noMethod         HandlersChain
	pool             sync.Pool
	trees            methodTrees
	maxParams        uint16
	maxSections      uint16
	trustedProxies   []string
	trustedCIDRs     []*net.IPNet
}

var _ IRouter = (*Engine)(nil)









func New(opts ...OptionFunc) *Engine {
	debugPrintWARNINGNew()
	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		FuncMap:                template.FuncMap{},
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      false,
		HandleMethodNotAllowed: false,
		ForwardedByClientIP:    true,
		RemoteIPHeaders:        []string{"X-Forwarded-For", "X-Real-IP"},
		TrustedPlatform:        defaultPlatform,
		UseRawPath:             false,
		RemoveExtraSlash:       false,
		UnescapePathValues:     true,
		MaxMultipartMemory:     defaultMultipartMemory,
		trees:                  make(methodTrees, 0, 9),
		delims:                 render.Delims{Left: "{{", Right: "}}"},
		secureJSONPrefix:       "while(1);",
		trustedProxies:         []string{"0.0.0.0/0", "::/0"},
		trustedCIDRs:           defaultTrustedCIDRs,
	}
	engine.RouterGroup.engine = engine
	engine.pool.New = func() any {
		return engine.allocateContext(engine.maxParams)
	}
	return engine.With(opts...)
}


func Default(opts ...OptionFunc) *Engine {
	debugPrintWARNINGDefault()
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine.With(opts...)
}

func (engine *Engine) Handler() http.Handler {
	if !engine.UseH2C {
		return engine
	}

	h2s := &http2.Server{}
	return h2c.NewHandler(engine, h2s)
}

func (engine *Engine) allocateContext(maxParams uint16) *Context {
	v := make(Params, 0, maxParams)
	skippedNodes := make([]skippedNode, 0, engine.maxSections)
	return &Context{engine: engine, params: &v, skippedNodes: &skippedNodes}
}


func (engine *Engine) Delims(left, right string) *Engine {
	engine.delims = render.Delims{Left: left, Right: right}
	return engine
}


func (engine *Engine) SecureJsonPrefix(prefix string) *Engine {
	engine.secureJSONPrefix = prefix
	return engine
}



func (engine *Engine) LoadHTMLGlob(pattern string) {
	left := engine.delims.Left
	right := engine.delims.Right
	templ := template.Must(template.New("").Delims(left, right).Funcs(engine.FuncMap).ParseGlob(pattern))

	if IsDebugging() {
		debugPrintLoadTemplate(templ)
		engine.HTMLRender = render.HTMLDebug{Glob: pattern, FuncMap: engine.FuncMap, Delims: engine.delims}
		return
	}

	engine.SetHTMLTemplate(templ)
}



func (engine *Engine) LoadHTMLFiles(files ...string) {
	if IsDebugging() {
		engine.HTMLRender = render.HTMLDebug{Files: files, FuncMap: engine.FuncMap, Delims: engine.delims}
		return
	}

	templ := template.Must(template.New("").Delims(engine.delims.Left, engine.delims.Right).Funcs(engine.FuncMap).ParseFiles(files...))
	engine.SetHTMLTemplate(templ)
}


func (engine *Engine) SetHTMLTemplate(templ *template.Template) {
	if len(engine.trees) > 0 {
		debugPrintWARNINGSetHTMLTemplate()
	}

	engine.HTMLRender = render.HTMLProduction{Template: templ.Funcs(engine.FuncMap)}
}


func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.FuncMap = funcMap
}


func (engine *Engine) NoRoute(handlers ...HandlerFunc) {
	engine.noRoute = handlers
	engine.rebuild404Handlers()
}


func (engine *Engine) NoMethod(handlers ...HandlerFunc) {
	engine.noMethod = handlers
	engine.rebuild405Handlers()
}




func (engine *Engine) Use(middleware ...HandlerFunc) IRoutes {
	engine.RouterGroup.Use(middleware...)
	engine.rebuild404Handlers()
	engine.rebuild405Handlers()
	return engine
}


func (engine *Engine) With(opts ...OptionFunc) *Engine {
	for _, opt := range opts {
		opt(engine)
	}

	return engine
}

func (engine *Engine) rebuild404Handlers() {
	engine.allNoRoute = engine.combineHandlers(engine.noRoute)
}

func (engine *Engine) rebuild405Handlers() {
	engine.allNoMethod = engine.combineHandlers(engine.noMethod)
}

func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
	assert1(path[0] == '/', "path must begin with '/'")
	assert1(method != "", "HTTP method can not be empty")
	assert1(len(handlers) > 0, "there must be at least one handler")

	debugPrintRoute(method, path, handlers)

	root := engine.trees.get(method)
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		engine.trees = append(engine.trees, methodTree{method: method, root: root})
	}
	root.addRoute(path, handlers)

	if paramsCount := countParams(path); paramsCount > engine.maxParams {
		engine.maxParams = paramsCount
	}

	if sectionsCount := countSections(path); sectionsCount > engine.maxSections {
		engine.maxSections = sectionsCount
	}
}



func (engine *Engine) Routes() (routes RoutesInfo) {
	for _, tree := range engine.trees {
		routes = iterate("", tree.method, routes, tree.root)
	}
	return routes
}

func iterate(path, method string, routes RoutesInfo, root *node) RoutesInfo {
	path += root.path
	if len(root.handlers) > 0 {
		handlerFunc := root.handlers.Last()
		routes = append(routes, RouteInfo{
			Method:      method,
			Path:        path,
			Handler:     nameOfFunction(handlerFunc),
			HandlerFunc: handlerFunc,
		})
	}
	for _, child := range root.children {
		routes = iterate(path, method, routes, child)
	}
	return routes
}




func (engine *Engine) Run(addr ...string) (err error) {
	defer func() { debugPrintError(err) }()

	if engine.isUnsafeTrustedProxies() {
		debugPrint("[WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.\n" +
			"Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.")
	}

	address := resolveAddress(addr)
	debugPrint("Listening and serving HTTP on %s\n", address)
	err = http.ListenAndServe(address, engine.Handler())
	return
}

func (engine *Engine) prepareTrustedCIDRs() ([]*net.IPNet, error) {
	if engine.trustedProxies == nil {
		return nil, nil
	}

	cidr := make([]*net.IPNet, 0, len(engine.trustedProxies))
	for _, trustedProxy := range engine.trustedProxies {
		if !strings.Contains(trustedProxy, "/") {
			ip := parseIP(trustedProxy)
			if ip == nil {
				return cidr, &net.ParseError{Type: "IP address", Text: trustedProxy}
			}

			switch len(ip) {
			case net.IPv4len:
				trustedProxy += "/32"
			case net.IPv6len:
				trustedProxy += "/128"
			}
		}
		_, cidrNet, err := net.ParseCIDR(trustedProxy)
		if err != nil {
			return cidr, err
		}
		cidr = append(cidr, cidrNet)
	}
	return cidr, nil
}









func (engine *Engine) SetTrustedProxies(trustedProxies []string) error {
	engine.trustedProxies = trustedProxies
	return engine.parseTrustedProxies()
}


func (engine *Engine) isUnsafeTrustedProxies() bool {
	return engine.isTrustedProxy(net.ParseIP("0.0.0.0")) || engine.isTrustedProxy(net.ParseIP("::"))
}


func (engine *Engine) parseTrustedProxies() error {
	trustedCIDRs, err := engine.prepareTrustedCIDRs()
	engine.trustedCIDRs = trustedCIDRs
	return err
}


func (engine *Engine) isTrustedProxy(ip net.IP) bool {
	if engine.trustedCIDRs == nil {
		return false
	}
	for _, cidr := range engine.trustedCIDRs {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}


func (engine *Engine) validateHeader(header string) (clientIP string, valid bool) {
	if header == "" {
		return "", false
	}
	items := strings.Split(header, ",")
	for i := len(items) - 1; i >= 0; i-- {
		ipStr := strings.TrimSpace(items[i])
		ip := net.ParseIP(ipStr)
		if ip == nil {
			break
		}

		
		
		if (i == 0) || (!engine.isTrustedProxy(ip)) {
			return ipStr, true
		}
	}
	return "", false
}



func parseIP(ip string) net.IP {
	parsedIP := net.ParseIP(ip)

	if ipv4 := parsedIP.To4(); ipv4 != nil {
		
		return ipv4
	}

	
	return parsedIP
}




func (engine *Engine) RunTLS(addr, certFile, keyFile string) (err error) {
	debugPrint("Listening and serving HTTPS on %s\n", addr)
	defer func() { debugPrintError(err) }()

	if engine.isUnsafeTrustedProxies() {
		debugPrint("[WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.\n" +
			"Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.")
	}

	err = http.ListenAndServeTLS(addr, certFile, keyFile, engine.Handler())
	return
}




func (engine *Engine) RunUnix(file string) (err error) {
	debugPrint("Listening and serving HTTP on unix:/%s", file)
	defer func() { debugPrintError(err) }()

	if engine.isUnsafeTrustedProxies() {
		debugPrint("[WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.\n" +
			"Please check https://github.com/gin-gonic/gin/blob/master/docs/doc.md#dont-trust-all-proxies for details.")
	}

	listener, err := net.Listen("unix", file)
	if err != nil {
		return
	}
	defer listener.Close()
	defer os.Remove(file)

	err = http.Serve(listener, engine.Handler())
	return
}




func (engine *Engine) RunFd(fd int) (err error) {
	debugPrint("Listening and serving HTTP on fd@%d", fd)
	defer func() { debugPrintError(err) }()

	if engine.isUnsafeTrustedProxies() {
		debugPrint("[WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.\n" +
			"Please check https://github.com/gin-gonic/gin/blob/master/docs/doc.md#dont-trust-all-proxies for details.")
	}

	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd@%d", fd))
	listener, err := net.FileListener(f)
	if err != nil {
		return
	}
	defer listener.Close()
	err = engine.RunListener(listener)
	return
}



func (engine *Engine) RunListener(listener net.Listener) (err error) {
	debugPrint("Listening and serving HTTP on listener what's bind with address@%s", listener.Addr())
	defer func() { debugPrintError(err) }()

	if engine.isUnsafeTrustedProxies() {
		debugPrint("[WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.\n" +
			"Please check https://github.com/gin-gonic/gin/blob/master/docs/doc.md#dont-trust-all-proxies for details.")
	}

	err = http.Serve(listener, engine.Handler())
	return
}


func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := engine.pool.Get().(*Context)
	c.writermem.reset(w)
	c.Request = req
	c.reset()

	engine.handleHTTPRequest(c)

	engine.pool.Put(c)
}




func (engine *Engine) HandleContext(c *Context) {
	oldIndexValue := c.index
	c.reset()
	engine.handleHTTPRequest(c)

	c.index = oldIndexValue
}

func (engine *Engine) handleHTTPRequest(c *Context) {
	httpMethod := c.Request.Method
	rPath := c.Request.URL.Path
	unescape := false
	if engine.UseRawPath && len(c.Request.URL.RawPath) > 0 {
		rPath = c.Request.URL.RawPath
		unescape = engine.UnescapePathValues
	}

	if engine.RemoveExtraSlash {
		rPath = cleanPath(rPath)
	}

	
	t := engine.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method != httpMethod {
			continue
		}
		root := t[i].root
		
		value := root.getValue(rPath, c.params, c.skippedNodes, unescape)
		if value.params != nil {
			c.Params = *value.params
		}
		if value.handlers != nil {
			c.handlers = value.handlers
			c.fullPath = value.fullPath
			c.Next()
			c.writermem.WriteHeaderNow()
			return
		}
		if httpMethod != http.MethodConnect && rPath != "/" {
			if value.tsr && engine.RedirectTrailingSlash {
				redirectTrailingSlash(c)
				return
			}
			if engine.RedirectFixedPath && redirectFixedPath(c, root, engine.RedirectFixedPath) {
				return
			}
		}
		break
	}

	if engine.HandleMethodNotAllowed {
		
		
		allowed := make([]string, 0, len(t)-1)
		for _, tree := range engine.trees {
			if tree.method == httpMethod {
				continue
			}
			if value := tree.root.getValue(rPath, nil, c.skippedNodes, unescape); value.handlers != nil {
				allowed = append(allowed, tree.method)
			}
		}
		if len(allowed) > 0 {
			c.handlers = engine.allNoMethod
			c.writermem.Header().Set("Allow", strings.Join(allowed, ", "))
			serveError(c, http.StatusMethodNotAllowed, default405Body)
			return
		}
	}

	c.handlers = engine.allNoRoute
	serveError(c, http.StatusNotFound, default404Body)
}

var mimePlain = []string{MIMEPlain}

func serveError(c *Context, code int, defaultMessage []byte) {
	c.writermem.status = code
	c.Next()
	if c.writermem.Written() {
		return
	}
	if c.writermem.Status() == code {
		c.writermem.Header()["Content-Type"] = mimePlain
		_, err := c.Writer.Write(defaultMessage)
		if err != nil {
			debugPrint("cannot write message to writer during serve error: %v", err)
		}
		return
	}
	c.writermem.WriteHeaderNow()
}

func redirectTrailingSlash(c *Context) {
	req := c.Request
	p := req.URL.Path
	if prefix := path.Clean(c.Request.Header.Get("X-Forwarded-Prefix")); prefix != "." {
		prefix = regSafePrefix.ReplaceAllString(prefix, "")
		prefix = regRemoveRepeatedChar.ReplaceAllString(prefix, "/")

		p = prefix + "/" + req.URL.Path
	}
	req.URL.Path = p + "/"
	if length := len(p); length > 1 && p[length-1] == '/' {
		req.URL.Path = p[:length-1]
	}
	redirectRequest(c)
}

func redirectFixedPath(c *Context, root *node, trailingSlash bool) bool {
	req := c.Request
	rPath := req.URL.Path

	if fixedPath, ok := root.findCaseInsensitivePath(cleanPath(rPath), trailingSlash); ok {
		req.URL.Path = bytesconv.BytesToString(fixedPath)
		redirectRequest(c)
		return true
	}
	return false
}

func redirectRequest(c *Context) {
	req := c.Request
	rPath := req.URL.Path
	rURL := req.URL.String()

	code := http.StatusMovedPermanently 
	if req.Method != http.MethodGet {
		code = http.StatusTemporaryRedirect
	}
	debugPrint("redirecting request %d: %s --> %s", code, rPath, rURL)
	http.Redirect(c.Writer, req, rURL, code)
	c.writermem.WriteHeaderNow()
}
