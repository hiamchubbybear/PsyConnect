























package fsnotify

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)


































































type Watcher struct {
	b backend

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Events chan Event

	
	Errors chan error
}


type Event struct {
	
	
	
	
	
	Name string

	
	
	
	
	Op Op

	
	
	
	
	
	
	
	
	renamedFrom string
}


type Op uint32



const (
	
	Create Op = 1 << iota

	
	
	Write

	
	
	
	Remove

	
	
	Rename

	
	
	
	
	
	Chmod

	
	
	
	xUnportableOpen

	
	
	
	xUnportableRead

	
	
	
	
	
	
	
	
	xUnportableCloseWrite

	
	
	
	xUnportableCloseRead
)

var (
	
	
	ErrNonExistentWatch = errors.New("fsnotify: can't remove non-existent watch")

	
	ErrClosed = errors.New("fsnotify: watcher already closed")

	
	
	
	
	
	
	
	
	ErrEventOverflow = errors.New("fsnotify: queue or buffer overflow")

	
	
	xErrUnsupported = errors.New("fsnotify: not supported with this backend")
)


func NewWatcher() (*Watcher, error) {
	ev, errs := make(chan Event), make(chan error)
	b, err := newBackend(ev, errs)
	if err != nil {
		return nil, err
	}
	return &Watcher{b: b, Events: ev, Errors: errs}, nil
}









func NewBufferedWatcher(sz uint) (*Watcher, error) {
	ev, errs := make(chan Event), make(chan error)
	b, err := newBufferedBackend(sz, ev, errs)
	if err != nil {
		return nil, err
	}
	return &Watcher{b: b, Events: ev, Errors: errs}, nil
}





































func (w *Watcher) Add(path string) error { return w.b.Add(path) }








func (w *Watcher) AddWith(path string, opts ...addOpt) error { return w.b.AddWith(path, opts...) }









func (w *Watcher) Remove(path string) error { return w.b.Remove(path) }


func (w *Watcher) Close() error { return w.b.Close() }





func (w *Watcher) WatchList() []string { return w.b.WatchList() }





func (w *Watcher) xSupports(op Op) bool { return w.b.xSupports(op) }

func (o Op) String() string {
	var b strings.Builder
	if o.Has(Create) {
		b.WriteString("|CREATE")
	}
	if o.Has(Remove) {
		b.WriteString("|REMOVE")
	}
	if o.Has(Write) {
		b.WriteString("|WRITE")
	}
	if o.Has(xUnportableOpen) {
		b.WriteString("|OPEN")
	}
	if o.Has(xUnportableRead) {
		b.WriteString("|READ")
	}
	if o.Has(xUnportableCloseWrite) {
		b.WriteString("|CLOSE_WRITE")
	}
	if o.Has(xUnportableCloseRead) {
		b.WriteString("|CLOSE_READ")
	}
	if o.Has(Rename) {
		b.WriteString("|RENAME")
	}
	if o.Has(Chmod) {
		b.WriteString("|CHMOD")
	}
	if b.Len() == 0 {
		return "[no events]"
	}
	return b.String()[1:]
}


func (o Op) Has(h Op) bool { return o&h != 0 }


func (e Event) Has(op Op) bool { return e.Op.Has(op) }


func (e Event) String() string {
	if e.renamedFrom != "" {
		return fmt.Sprintf("%-13s %q ← %q", e.Op.String(), e.Name, e.renamedFrom)
	}
	return fmt.Sprintf("%-13s %q", e.Op.String(), e.Name)
}

type (
	backend interface {
		Add(string) error
		AddWith(string, ...addOpt) error
		Remove(string) error
		WatchList() []string
		Close() error
		xSupports(Op) bool
	}
	addOpt   func(opt *withOpts)
	withOpts struct {
		bufsize    int
		op         Op
		noFollow   bool
		sendCreate bool
	}
)

var debug = func() bool {
	
	
	
	return os.Getenv("FSNOTIFY_DEBUG") == "1"
}()

var defaultOpts = withOpts{
	bufsize: 65536, 
	op:      Create | Write | Remove | Rename | Chmod,
}

func getOptions(opts ...addOpt) withOpts {
	with := defaultOpts
	for _, o := range opts {
		if o != nil {
			o(&with)
		}
	}
	return with
}











func WithBufferSize(bytes int) addOpt {
	return func(opt *withOpts) { opt.bufsize = bytes }
}















func withOps(op Op) addOpt {
	return func(opt *withOpts) { opt.op = op }
}



func withNoFollow() addOpt {
	return func(opt *withOpts) { opt.noFollow = true }
}


func withCreate() addOpt {
	return func(opt *withOpts) { opt.sendCreate = true }
}

var enableRecurse = false



func recursivePath(path string) (string, bool) {
	path = filepath.Clean(path)
	if !enableRecurse { 
		return path, false
	}
	if filepath.Base(path) == "..." {
		return filepath.Dir(path), true
	}
	return path, false
}
