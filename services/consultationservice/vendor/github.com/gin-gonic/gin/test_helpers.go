



package gin

import "net/http"


func CreateTestContext(w http.ResponseWriter) (c *Context, r *Engine) {
	r = New()
	c = r.allocateContext(0)
	c.reset()
	c.writermem.reset(w)
	return
}


func CreateTestContextOnly(w http.ResponseWriter, r *Engine) (c *Context) {
	c = r.allocateContext(r.maxParams)
	c.reset()
	c.writermem.reset(w)
	return
}
