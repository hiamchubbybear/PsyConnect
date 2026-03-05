

package grpclog

import (
	"fmt"
)


type componentData struct {
	name string
}

var cache = map[string]*componentData{}

func (c *componentData) InfoDepth(depth int, args ...any) {
	args = append([]any{"[" + string(c.name) + "]"}, args...)
	InfoDepth(depth+1, args...)
}

func (c *componentData) WarningDepth(depth int, args ...any) {
	args = append([]any{"[" + string(c.name) + "]"}, args...)
	WarningDepth(depth+1, args...)
}

func (c *componentData) ErrorDepth(depth int, args ...any) {
	args = append([]any{"[" + string(c.name) + "]"}, args...)
	ErrorDepth(depth+1, args...)
}

func (c *componentData) FatalDepth(depth int, args ...any) {
	args = append([]any{"[" + string(c.name) + "]"}, args...)
	FatalDepth(depth+1, args...)
}

func (c *componentData) Info(args ...any) {
	c.InfoDepth(1, args...)
}

func (c *componentData) Warning(args ...any) {
	c.WarningDepth(1, args...)
}

func (c *componentData) Error(args ...any) {
	c.ErrorDepth(1, args...)
}

func (c *componentData) Fatal(args ...any) {
	c.FatalDepth(1, args...)
}

func (c *componentData) Infof(format string, args ...any) {
	c.InfoDepth(1, fmt.Sprintf(format, args...))
}

func (c *componentData) Warningf(format string, args ...any) {
	c.WarningDepth(1, fmt.Sprintf(format, args...))
}

func (c *componentData) Errorf(format string, args ...any) {
	c.ErrorDepth(1, fmt.Sprintf(format, args...))
}

func (c *componentData) Fatalf(format string, args ...any) {
	c.FatalDepth(1, fmt.Sprintf(format, args...))
}

func (c *componentData) Infoln(args ...any) {
	c.InfoDepth(1, args...)
}

func (c *componentData) Warningln(args ...any) {
	c.WarningDepth(1, args...)
}

func (c *componentData) Errorln(args ...any) {
	c.ErrorDepth(1, args...)
}

func (c *componentData) Fatalln(args ...any) {
	c.FatalDepth(1, args...)
}

func (c *componentData) V(l int) bool {
	return V(l)
}





func Component(componentName string) DepthLoggerV2 {
	if cData, ok := cache[componentName]; ok {
		return cData
	}
	c := &componentData{componentName}
	cache[componentName] = c
	return c
}
