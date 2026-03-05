//go:build linux && !appengine

package fsnotify

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/fsnotify/fsnotify/internal"
	"golang.org/x/sys/unix"
)

type inotify struct {
	Events chan Event
	Errors chan error

	
	
	fd          int
	inotifyFile *os.File
	watches     *watches
	done        chan struct{} 
	doneMu      sync.Mutex
	doneResp    chan struct{} 

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	cookies     [10]koekje
	cookieIndex uint8
	cookiesMu   sync.Mutex
}

type (
	watches struct {
		mu   sync.RWMutex
		wd   map[uint32]*watch 
		path map[string]uint32 
	}
	watch struct {
		wd      uint32 
		flags   uint32 
		path    string 
		recurse bool   
	}
	koekje struct {
		cookie uint32
		path   string
	}
)

func newWatches() *watches {
	return &watches{
		wd:   make(map[uint32]*watch),
		path: make(map[string]uint32),
	}
}

func (w *watches) len() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.wd)
}

func (w *watches) add(ww *watch) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.wd[ww.wd] = ww
	w.path[ww.path] = ww.wd
}

func (w *watches) remove(wd uint32) {
	w.mu.Lock()
	defer w.mu.Unlock()
	watch := w.wd[wd] 
	if watch == nil {
		return
	}
	delete(w.path, watch.path)
	delete(w.wd, wd)
}

func (w *watches) removePath(path string) ([]uint32, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	path, recurse := recursivePath(path)
	wd, ok := w.path[path]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNonExistentWatch, path)
	}

	watch := w.wd[wd]
	if recurse && !watch.recurse {
		return nil, fmt.Errorf("can't use /... with non-recursive watch %q", path)
	}

	delete(w.path, path)
	delete(w.wd, wd)
	if !watch.recurse {
		return []uint32{wd}, nil
	}

	wds := make([]uint32, 0, 8)
	wds = append(wds, wd)
	for p, rwd := range w.path {
		if filepath.HasPrefix(p, path) {
			delete(w.path, p)
			delete(w.wd, rwd)
			wds = append(wds, rwd)
		}
	}
	return wds, nil
}

func (w *watches) byPath(path string) *watch {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.wd[w.path[path]]
}

func (w *watches) byWd(wd uint32) *watch {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.wd[wd]
}

func (w *watches) updatePath(path string, f func(*watch) (*watch, error)) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	var existing *watch
	wd, ok := w.path[path]
	if ok {
		existing = w.wd[wd]
	}

	upd, err := f(existing)
	if err != nil {
		return err
	}
	if upd != nil {
		w.wd[upd.wd] = upd
		w.path[upd.path] = upd.wd

		if upd.wd != wd {
			delete(w.wd, wd)
		}
	}

	return nil
}

func newBackend(ev chan Event, errs chan error) (backend, error) {
	return newBufferedBackend(0, ev, errs)
}

func newBufferedBackend(sz uint, ev chan Event, errs chan error) (backend, error) {
	
	
	fd, errno := unix.InotifyInit1(unix.IN_CLOEXEC | unix.IN_NONBLOCK)
	if fd == -1 {
		return nil, errno
	}

	w := &inotify{
		Events:      ev,
		Errors:      errs,
		fd:          fd,
		inotifyFile: os.NewFile(uintptr(fd), ""),
		watches:     newWatches(),
		done:        make(chan struct{}),
		doneResp:    make(chan struct{}),
	}

	go w.readEvents()
	return w, nil
}


func (w *inotify) sendEvent(e Event) bool {
	select {
	case <-w.done:
		return false
	case w.Events <- e:
		return true
	}
}


func (w *inotify) sendError(err error) bool {
	if err == nil {
		return true
	}
	select {
	case <-w.done:
		return false
	case w.Errors <- err:
		return true
	}
}

func (w *inotify) isClosed() bool {
	select {
	case <-w.done:
		return true
	default:
		return false
	}
}

func (w *inotify) Close() error {
	w.doneMu.Lock()
	if w.isClosed() {
		w.doneMu.Unlock()
		return nil
	}
	close(w.done)
	w.doneMu.Unlock()

	
	
	err := w.inotifyFile.Close()
	if err != nil {
		return err
	}

	
	<-w.doneResp

	return nil
}

func (w *inotify) Add(name string) error { return w.AddWith(name) }

func (w *inotify) AddWith(path string, opts ...addOpt) error {
	if w.isClosed() {
		return ErrClosed
	}
	if debug {
		fmt.Fprintf(os.Stderr, "FSNOTIFY_DEBUG: %s  AddWith(%q)\n",
			time.Now().Format("15:04:05.000000000"), path)
	}

	with := getOptions(opts...)
	if !w.xSupports(with.op) {
		return fmt.Errorf("%w: %s", xErrUnsupported, with.op)
	}

	path, recurse := recursivePath(path)
	if recurse {
		return filepath.WalkDir(path, func(root string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				if root == path {
					return fmt.Errorf("fsnotify: not a directory: %q", path)
				}
				return nil
			}

			
			
			
			
			
			
			if with.sendCreate && root != path {
				w.sendEvent(Event{Name: root, Op: Create})
			}

			return w.add(root, with, true)
		})
	}

	return w.add(path, with, false)
}

func (w *inotify) add(path string, with withOpts, recurse bool) error {
	var flags uint32
	if with.noFollow {
		flags |= unix.IN_DONT_FOLLOW
	}
	if with.op.Has(Create) {
		flags |= unix.IN_CREATE
	}
	if with.op.Has(Write) {
		flags |= unix.IN_MODIFY
	}
	if with.op.Has(Remove) {
		flags |= unix.IN_DELETE | unix.IN_DELETE_SELF
	}
	if with.op.Has(Rename) {
		flags |= unix.IN_MOVED_TO | unix.IN_MOVED_FROM | unix.IN_MOVE_SELF
	}
	if with.op.Has(Chmod) {
		flags |= unix.IN_ATTRIB
	}
	if with.op.Has(xUnportableOpen) {
		flags |= unix.IN_OPEN
	}
	if with.op.Has(xUnportableRead) {
		flags |= unix.IN_ACCESS
	}
	if with.op.Has(xUnportableCloseWrite) {
		flags |= unix.IN_CLOSE_WRITE
	}
	if with.op.Has(xUnportableCloseRead) {
		flags |= unix.IN_CLOSE_NOWRITE
	}
	return w.register(path, flags, recurse)
}

func (w *inotify) register(path string, flags uint32, recurse bool) error {
	return w.watches.updatePath(path, func(existing *watch) (*watch, error) {
		if existing != nil {
			flags |= existing.flags | unix.IN_MASK_ADD
		}

		wd, err := unix.InotifyAddWatch(w.fd, path, flags)
		if wd == -1 {
			return nil, err
		}

		if existing == nil {
			return &watch{
				wd:      uint32(wd),
				path:    path,
				flags:   flags,
				recurse: recurse,
			}, nil
		}

		existing.wd = uint32(wd)
		existing.flags = flags
		return existing, nil
	})
}

func (w *inotify) Remove(name string) error {
	if w.isClosed() {
		return nil
	}
	if debug {
		fmt.Fprintf(os.Stderr, "FSNOTIFY_DEBUG: %s  Remove(%q)\n",
			time.Now().Format("15:04:05.000000000"), name)
	}
	return w.remove(filepath.Clean(name))
}

func (w *inotify) remove(name string) error {
	wds, err := w.watches.removePath(name)
	if err != nil {
		return err
	}

	for _, wd := range wds {
		_, err := unix.InotifyRmWatch(w.fd, wd)
		if err != nil {
			
			
			
			
			
			
			
			
			
			
			
			return err
		}
	}
	return nil
}

func (w *inotify) WatchList() []string {
	if w.isClosed() {
		return nil
	}

	entries := make([]string, 0, w.watches.len())
	w.watches.mu.RLock()
	for pathname := range w.watches.path {
		entries = append(entries, pathname)
	}
	w.watches.mu.RUnlock()

	return entries
}



func (w *inotify) readEvents() {
	defer func() {
		close(w.doneResp)
		close(w.Errors)
		close(w.Events)
	}()

	var (
		buf   [unix.SizeofInotifyEvent * 4096]byte 
		errno error                                
	)
	for {
		
		if w.isClosed() {
			return
		}

		n, err := w.inotifyFile.Read(buf[:])
		switch {
		case errors.Unwrap(err) == os.ErrClosed:
			return
		case err != nil:
			if !w.sendError(err) {
				return
			}
			continue
		}

		if n < unix.SizeofInotifyEvent {
			var err error
			if n == 0 {
				err = io.EOF 
			} else if n < 0 {
				err = errno 
			} else {
				err = errors.New("notify: short read in readEvents()") 
			}
			if !w.sendError(err) {
				return
			}
			continue
		}

		
		
		var offset uint32
		for offset <= uint32(n-unix.SizeofInotifyEvent) {
			var (
				
				raw     = (*unix.InotifyEvent)(unsafe.Pointer(&buf[offset]))
				mask    = uint32(raw.Mask)
				nameLen = uint32(raw.Len)
				
				next = func() { offset += unix.SizeofInotifyEvent + nameLen }
			)

			if mask&unix.IN_Q_OVERFLOW != 0 {
				if !w.sendError(ErrEventOverflow) {
					return
				}
			}

			
			
			
			
			
			watch := w.watches.byWd(uint32(raw.Wd))
			
			
			
			
			if watch == nil {
				next()
				continue
			}

			name := watch.path
			if nameLen > 0 {
				
				bytes := (*[unix.PathMax]byte)(unsafe.Pointer(&buf[offset+unix.SizeofInotifyEvent]))[:nameLen:nameLen]
				
				name += "/" + strings.TrimRight(string(bytes[0:nameLen]), "\000")
			}

			if debug {
				internal.Debug(name, raw.Mask, raw.Cookie)
			}

			if mask&unix.IN_IGNORED != 0 { 
				next()
				continue
			}

			
			
			if mask&unix.IN_DELETE_SELF == unix.IN_DELETE_SELF {
				w.watches.remove(watch.wd)
			}

			
			
			
			if mask&unix.IN_MOVE_SELF == unix.IN_MOVE_SELF {
				if watch.recurse {
					next() 
					continue
				}

				err := w.remove(watch.path)
				if err != nil && !errors.Is(err, ErrNonExistentWatch) {
					if !w.sendError(err) {
						return
					}
				}
			}

			
			
			if mask&unix.IN_DELETE_SELF != 0 {
				if _, ok := w.watches.path[filepath.Dir(watch.path)]; ok {
					next()
					continue
				}
			}

			ev := w.newEvent(name, mask, raw.Cookie)
			
			if watch.recurse {
				isDir := mask&unix.IN_ISDIR == unix.IN_ISDIR
				
				if isDir && ev.Has(Create) {
					err := w.register(ev.Name, watch.flags, true)
					if !w.sendError(err) {
						return
					}

					
					
					
					
					
					
					
					
					
					if ev.renamedFrom != "" {
						w.watches.mu.Lock()
						for k, ww := range w.watches.wd {
							if k == watch.wd || ww.path == ev.Name {
								continue
							}
							if strings.HasPrefix(ww.path, ev.renamedFrom) {
								ww.path = strings.Replace(ww.path, ev.renamedFrom, ev.Name, 1)
								w.watches.wd[k] = ww
							}
						}
						w.watches.mu.Unlock()
					}
				}
			}

			
			if !w.sendEvent(ev) {
				return
			}
			next()
		}
	}
}

func (w *inotify) isRecursive(path string) bool {
	ww := w.watches.byPath(path)
	if ww == nil { 
		ww = w.watches.byPath(filepath.Dir(path))
	}
	return ww != nil && ww.recurse
}

func (w *inotify) newEvent(name string, mask, cookie uint32) Event {
	e := Event{Name: name}
	if mask&unix.IN_CREATE == unix.IN_CREATE || mask&unix.IN_MOVED_TO == unix.IN_MOVED_TO {
		e.Op |= Create
	}
	if mask&unix.IN_DELETE_SELF == unix.IN_DELETE_SELF || mask&unix.IN_DELETE == unix.IN_DELETE {
		e.Op |= Remove
	}
	if mask&unix.IN_MODIFY == unix.IN_MODIFY {
		e.Op |= Write
	}
	if mask&unix.IN_OPEN == unix.IN_OPEN {
		e.Op |= xUnportableOpen
	}
	if mask&unix.IN_ACCESS == unix.IN_ACCESS {
		e.Op |= xUnportableRead
	}
	if mask&unix.IN_CLOSE_WRITE == unix.IN_CLOSE_WRITE {
		e.Op |= xUnportableCloseWrite
	}
	if mask&unix.IN_CLOSE_NOWRITE == unix.IN_CLOSE_NOWRITE {
		e.Op |= xUnportableCloseRead
	}
	if mask&unix.IN_MOVE_SELF == unix.IN_MOVE_SELF || mask&unix.IN_MOVED_FROM == unix.IN_MOVED_FROM {
		e.Op |= Rename
	}
	if mask&unix.IN_ATTRIB == unix.IN_ATTRIB {
		e.Op |= Chmod
	}

	if cookie != 0 {
		if mask&unix.IN_MOVED_FROM == unix.IN_MOVED_FROM {
			w.cookiesMu.Lock()
			w.cookies[w.cookieIndex] = koekje{cookie: cookie, path: e.Name}
			w.cookieIndex++
			if w.cookieIndex > 9 {
				w.cookieIndex = 0
			}
			w.cookiesMu.Unlock()
		} else if mask&unix.IN_MOVED_TO == unix.IN_MOVED_TO {
			w.cookiesMu.Lock()
			var prev string
			for _, c := range w.cookies {
				if c.cookie == cookie {
					prev = c.path
					break
				}
			}
			w.cookiesMu.Unlock()
			e.renamedFrom = prev
		}
	}
	return e
}

func (w *inotify) xSupports(op Op) bool {
	return true 
}

func (w *inotify) state() {
	w.watches.mu.Lock()
	defer w.watches.mu.Unlock()
	for wd, ww := range w.watches.wd {
		fmt.Fprintf(os.Stderr, "%4d: recurse=%t %q\n", wd, ww.recurse, ww.path)
	}
}
