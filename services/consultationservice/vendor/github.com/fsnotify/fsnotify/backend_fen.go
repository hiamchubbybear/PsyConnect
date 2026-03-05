//go:build solaris





package fsnotify

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify/internal"
	"golang.org/x/sys/unix"
)

type fen struct {
	Events chan Event
	Errors chan error

	mu      sync.Mutex
	port    *unix.EventPort
	done    chan struct{} 
	dirs    map[string]Op 
	watches map[string]Op 
}

func newBackend(ev chan Event, errs chan error) (backend, error) {
	return newBufferedBackend(0, ev, errs)
}

func newBufferedBackend(sz uint, ev chan Event, errs chan error) (backend, error) {
	w := &fen{
		Events:  ev,
		Errors:  errs,
		dirs:    make(map[string]Op),
		watches: make(map[string]Op),
		done:    make(chan struct{}),
	}

	var err error
	w.port, err = unix.NewEventPort()
	if err != nil {
		return nil, fmt.Errorf("fsnotify.NewWatcher: %w", err)
	}

	go w.readEvents()
	return w, nil
}



func (w *fen) sendEvent(name string, op Op) (sent bool) {
	select {
	case <-w.done:
		return false
	case w.Events <- Event{Name: name, Op: op}:
		return true
	}
}



func (w *fen) sendError(err error) (sent bool) {
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

func (w *fen) isClosed() bool {
	select {
	case <-w.done:
		return true
	default:
		return false
	}
}

func (w *fen) Close() error {
	
	
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.isClosed() {
		return nil
	}
	close(w.done)
	return w.port.Close()
}

func (w *fen) Add(name string) error { return w.AddWith(name) }

func (w *fen) AddWith(name string, opts ...addOpt) error {
	if w.isClosed() {
		return ErrClosed
	}
	if debug {
		fmt.Fprintf(os.Stderr, "FSNOTIFY_DEBUG: %s  AddWith(%q)\n",
			time.Now().Format("15:04:05.000000000"), name)
	}

	with := getOptions(opts...)
	if !w.xSupports(with.op) {
		return fmt.Errorf("%w: %s", xErrUnsupported, with.op)
	}

	
	
	stat, err := os.Stat(name)
	if err != nil {
		return err
	}

	
	if stat.IsDir() {
		err := w.handleDirectory(name, stat, true, w.associateFile)
		if err != nil {
			return err
		}

		w.mu.Lock()
		w.dirs[name] = with.op
		w.mu.Unlock()
		return nil
	}

	err = w.associateFile(name, stat, true)
	if err != nil {
		return err
	}

	w.mu.Lock()
	w.watches[name] = with.op
	w.mu.Unlock()
	return nil
}

func (w *fen) Remove(name string) error {
	if w.isClosed() {
		return nil
	}
	if !w.port.PathIsWatched(name) {
		return fmt.Errorf("%w: %s", ErrNonExistentWatch, name)
	}
	if debug {
		fmt.Fprintf(os.Stderr, "FSNOTIFY_DEBUG: %s  Remove(%q)\n",
			time.Now().Format("15:04:05.000000000"), name)
	}

	
	
	
	w.mu.Lock()
	delete(w.watches, name)
	delete(w.dirs, name)
	w.mu.Unlock()

	stat, err := os.Stat(name)
	if err != nil {
		return err
	}

	
	if stat.IsDir() {
		err := w.handleDirectory(name, stat, false, w.dissociateFile)
		if err != nil {
			return err
		}
		return nil
	}

	err = w.port.DissociatePath(name)
	if err != nil {
		return err
	}

	return nil
}


func (w *fen) readEvents() {
	
	
	defer func() {
		close(w.Errors)
		close(w.Events)
	}()

	pevents := make([]unix.PortEvent, 8)
	for {
		count, err := w.port.Get(pevents, 1, nil)
		if err != nil && err != unix.ETIME {
			
			if errors.Is(err, unix.EINTR) && count == 0 {
				continue
			}
			
			if errors.Is(err, unix.EBADF) && w.isClosed() {
				return
			}
			
			if !w.sendError(err) {
				return
			}
		}

		p := pevents[:count]
		for _, pevent := range p {
			if pevent.Source != unix.PORT_SOURCE_FILE {
				
				if !w.sendError(errors.New("Event from unexpected source received")) {
					return
				}
				continue
			}

			if debug {
				internal.Debug(pevent.Path, pevent.Events)
			}

			err = w.handleEvent(&pevent)
			if !w.sendError(err) {
				return
			}
		}
	}
}

func (w *fen) handleDirectory(path string, stat os.FileInfo, follow bool, handler func(string, os.FileInfo, bool) error) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	
	for _, entry := range files {
		finfo, err := entry.Info()
		if err != nil {
			return err
		}
		err = handler(filepath.Join(path, finfo.Name()), finfo, false)
		if err != nil {
			return err
		}
	}

	
	return handler(path, stat, follow)
}





func (w *fen) handleEvent(event *unix.PortEvent) error {
	var (
		events     = event.Events
		path       = event.Path
		fmode      = event.Cookie.(os.FileMode)
		reRegister = true
	)

	w.mu.Lock()
	_, watchedDir := w.dirs[path]
	_, watchedPath := w.watches[path]
	w.mu.Unlock()
	isWatched := watchedDir || watchedPath

	if events&unix.FILE_DELETE != 0 {
		if !w.sendEvent(path, Remove) {
			return nil
		}
		reRegister = false
	}
	if events&unix.FILE_RENAME_FROM != 0 {
		if !w.sendEvent(path, Rename) {
			return nil
		}
		
		reRegister = false
	}
	if events&unix.FILE_RENAME_TO != 0 {
		
		
		
		

		
		
		if !w.sendEvent(path, Remove) {
			return nil
		}
		
		reRegister = false
	}

	
	if !reRegister {
		if watchedDir {
			w.mu.Lock()
			delete(w.dirs, path)
			w.mu.Unlock()
		}
		if watchedPath {
			w.mu.Lock()
			delete(w.watches, path)
			w.mu.Unlock()
		}
		return nil
	}

	
	
	

	stat, err := os.Lstat(path)
	if err != nil {
		
		
		
		
		
		
		if !w.sendEvent(path, Remove) {
			return nil
		}
		
		
		return nil
	}

	
	
	if isWatched {
		stat, err = os.Stat(path)
		if err != nil {
			
			
			if !w.sendEvent(path, Remove) {
				return nil
			}
			
		}
	}

	if events&unix.FILE_MODIFIED != 0 {
		if fmode.IsDir() && watchedDir {
			if err := w.updateDirectory(path); err != nil {
				return err
			}
		} else {
			if !w.sendEvent(path, Write) {
				return nil
			}
		}
	}
	if events&unix.FILE_ATTRIB != 0 && stat != nil {
		
		if stat.Mode().Perm() != fmode.Perm() {
			if !w.sendEvent(path, Chmod) {
				return nil
			}
		}
	}

	if stat != nil {
		
		
		return w.associateFile(path, stat, isWatched)
	}
	return nil
}

func (w *fen) updateDirectory(path string) error {
	
	
	
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range files {
		path := filepath.Join(path, entry.Name())
		if w.port.PathIsWatched(path) {
			continue
		}

		finfo, err := entry.Info()
		if err != nil {
			return err
		}
		err = w.associateFile(path, finfo, false)
		if !w.sendError(err) {
			return nil
		}
		if !w.sendEvent(path, Create) {
			return nil
		}
	}
	return nil
}

func (w *fen) associateFile(path string, stat os.FileInfo, follow bool) error {
	if w.isClosed() {
		return ErrClosed
	}
	
	
	
	
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.port.PathIsWatched(path) {
		
		
		
		
		
		err := w.port.DissociatePath(path)
		if err != nil && !errors.Is(err, unix.ENOENT) {
			return err
		}
	}

	var events int
	if !follow {
		
		
		events |= unix.FILE_NOFOLLOW
	}
	if true { 
		events |= unix.FILE_MODIFIED
	}
	if true {
		events |= unix.FILE_ATTRIB
	}
	return w.port.AssociatePath(path, stat, events, stat.Mode())
}

func (w *fen) dissociateFile(path string, stat os.FileInfo, unused bool) error {
	if !w.port.PathIsWatched(path) {
		return nil
	}
	return w.port.DissociatePath(path)
}

func (w *fen) WatchList() []string {
	if w.isClosed() {
		return nil
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	entries := make([]string, 0, len(w.watches)+len(w.dirs))
	for pathname := range w.dirs {
		entries = append(entries, pathname)
	}
	for pathname := range w.watches {
		entries = append(entries, pathname)
	}

	return entries
}

func (w *fen) xSupports(op Op) bool {
	if op.Has(xUnportableOpen) || op.Has(xUnportableRead) ||
		op.Has(xUnportableCloseWrite) || op.Has(xUnportableCloseRead) {
		return false
	}
	return true
}
