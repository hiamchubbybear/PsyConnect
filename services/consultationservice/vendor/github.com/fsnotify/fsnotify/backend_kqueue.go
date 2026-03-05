//go:build freebsd || openbsd || netbsd || dragonfly || darwin

package fsnotify

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify/internal"
	"golang.org/x/sys/unix"
)

type kqueue struct {
	Events chan Event
	Errors chan error

	kq        int    
	closepipe [2]int 
	watches   *watches
	done      chan struct{}
	doneMu    sync.Mutex
}

type (
	watches struct {
		mu     sync.RWMutex
		wd     map[int]watch               
		path   map[string]int              
		byDir  map[string]map[int]struct{} 
		seen   map[string]struct{}         
		byUser map[string]struct{}         
	}
	watch struct {
		wd       int
		name     string
		linkName string 
		isDir    bool
		dirFlags uint32
	}
)

func newWatches() *watches {
	return &watches{
		wd:     make(map[int]watch),
		path:   make(map[string]int),
		byDir:  make(map[string]map[int]struct{}),
		seen:   make(map[string]struct{}),
		byUser: make(map[string]struct{}),
	}
}

func (w *watches) listPaths(userOnly bool) []string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if userOnly {
		l := make([]string, 0, len(w.byUser))
		for p := range w.byUser {
			l = append(l, p)
		}
		return l
	}

	l := make([]string, 0, len(w.path))
	for p := range w.path {
		l = append(l, p)
	}
	return l
}

func (w *watches) watchesInDir(path string) []string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	l := make([]string, 0, 4)
	for fd := range w.byDir[path] {
		info := w.wd[fd]
		if _, ok := w.byUser[info.name]; !ok {
			l = append(l, info.name)
		}
	}
	return l
}


func (w *watches) addUserWatch(path string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.byUser[path] = struct{}{}
}

func (w *watches) addLink(path string, fd int) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.path[path] = fd
	w.seen[path] = struct{}{}
}

func (w *watches) add(path, linkPath string, fd int, isDir bool) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.path[path] = fd
	w.wd[fd] = watch{wd: fd, name: path, linkName: linkPath, isDir: isDir}

	parent := filepath.Dir(path)
	byDir, ok := w.byDir[parent]
	if !ok {
		byDir = make(map[int]struct{}, 1)
		w.byDir[parent] = byDir
	}
	byDir[fd] = struct{}{}
}

func (w *watches) byWd(fd int) (watch, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	info, ok := w.wd[fd]
	return info, ok
}

func (w *watches) byPath(path string) (watch, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	info, ok := w.wd[w.path[path]]
	return info, ok
}

func (w *watches) updateDirFlags(path string, flags uint32) {
	w.mu.Lock()
	defer w.mu.Unlock()

	fd := w.path[path]
	info := w.wd[fd]
	info.dirFlags = flags
	w.wd[fd] = info
}

func (w *watches) remove(fd int, path string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	isDir := w.wd[fd].isDir
	delete(w.path, path)
	delete(w.byUser, path)

	parent := filepath.Dir(path)
	delete(w.byDir[parent], fd)

	if len(w.byDir[parent]) == 0 {
		delete(w.byDir, parent)
	}

	delete(w.wd, fd)
	delete(w.seen, path)
	return isDir
}

func (w *watches) markSeen(path string, exists bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if exists {
		w.seen[path] = struct{}{}
	} else {
		delete(w.seen, path)
	}
}

func (w *watches) seenBefore(path string) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	_, ok := w.seen[path]
	return ok
}

func newBackend(ev chan Event, errs chan error) (backend, error) {
	return newBufferedBackend(0, ev, errs)
}

func newBufferedBackend(sz uint, ev chan Event, errs chan error) (backend, error) {
	kq, closepipe, err := newKqueue()
	if err != nil {
		return nil, err
	}

	w := &kqueue{
		Events:    ev,
		Errors:    errs,
		kq:        kq,
		closepipe: closepipe,
		done:      make(chan struct{}),
		watches:   newWatches(),
	}

	go w.readEvents()
	return w, nil
}







func newKqueue() (kq int, closepipe [2]int, err error) {
	kq, err = unix.Kqueue()
	if kq == -1 {
		return kq, closepipe, err
	}

	
	err = unix.Pipe(closepipe[:])
	if err != nil {
		unix.Close(kq)
		return kq, closepipe, err
	}
	unix.CloseOnExec(closepipe[0])
	unix.CloseOnExec(closepipe[1])

	
	changes := make([]unix.Kevent_t, 1)
	
	unix.SetKevent(&changes[0], closepipe[0], unix.EVFILT_READ,
		unix.EV_ADD|unix.EV_ENABLE|unix.EV_ONESHOT)

	ok, err := unix.Kevent(kq, changes, nil, nil)
	if ok == -1 {
		unix.Close(kq)
		unix.Close(closepipe[0])
		unix.Close(closepipe[1])
		return kq, closepipe, err
	}
	return kq, closepipe, nil
}


func (w *kqueue) sendEvent(e Event) bool {
	select {
	case <-w.done:
		return false
	case w.Events <- e:
		return true
	}
}


func (w *kqueue) sendError(err error) bool {
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

func (w *kqueue) isClosed() bool {
	select {
	case <-w.done:
		return true
	default:
		return false
	}
}

func (w *kqueue) Close() error {
	w.doneMu.Lock()
	if w.isClosed() {
		w.doneMu.Unlock()
		return nil
	}
	close(w.done)
	w.doneMu.Unlock()

	pathsToRemove := w.watches.listPaths(false)
	for _, name := range pathsToRemove {
		w.Remove(name)
	}

	
	unix.Close(w.closepipe[1])
	return nil
}

func (w *kqueue) Add(name string) error { return w.AddWith(name) }

func (w *kqueue) AddWith(name string, opts ...addOpt) error {
	if debug {
		fmt.Fprintf(os.Stderr, "FSNOTIFY_DEBUG: %s  AddWith(%q)\n",
			time.Now().Format("15:04:05.000000000"), name)
	}

	with := getOptions(opts...)
	if !w.xSupports(with.op) {
		return fmt.Errorf("%w: %s", xErrUnsupported, with.op)
	}

	_, err := w.addWatch(name, noteAllEvents)
	if err != nil {
		return err
	}
	w.watches.addUserWatch(name)
	return nil
}

func (w *kqueue) Remove(name string) error {
	if debug {
		fmt.Fprintf(os.Stderr, "FSNOTIFY_DEBUG: %s  Remove(%q)\n",
			time.Now().Format("15:04:05.000000000"), name)
	}
	return w.remove(name, true)
}

func (w *kqueue) remove(name string, unwatchFiles bool) error {
	if w.isClosed() {
		return nil
	}

	name = filepath.Clean(name)
	info, ok := w.watches.byPath(name)
	if !ok {
		return fmt.Errorf("%w: %s", ErrNonExistentWatch, name)
	}

	err := w.register([]int{info.wd}, unix.EV_DELETE, 0)
	if err != nil {
		return err
	}

	unix.Close(info.wd)

	isDir := w.watches.remove(info.wd, name)

	
	if unwatchFiles && isDir {
		pathsToRemove := w.watches.watchesInDir(name)
		for _, name := range pathsToRemove {
			
			
			
			w.Remove(name)
		}
	}
	return nil
}

func (w *kqueue) WatchList() []string {
	if w.isClosed() {
		return nil
	}
	return w.watches.listPaths(true)
}


const noteAllEvents = unix.NOTE_DELETE | unix.NOTE_WRITE | unix.NOTE_ATTRIB | unix.NOTE_RENAME





func (w *kqueue) addWatch(name string, flags uint32) (string, error) {
	if w.isClosed() {
		return "", ErrClosed
	}

	name = filepath.Clean(name)

	info, alreadyWatching := w.watches.byPath(name)
	if !alreadyWatching {
		fi, err := os.Lstat(name)
		if err != nil {
			return "", err
		}

		
		if (fi.Mode()&os.ModeSocket == os.ModeSocket) || (fi.Mode()&os.ModeNamedPipe == os.ModeNamedPipe) {
			return "", nil
		}

		
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			link, err := os.Readlink(name)
			if err != nil {
				
				
				
				
				return "", nil
			}

			_, alreadyWatching = w.watches.byPath(link)
			if alreadyWatching {
				
				
				w.watches.addLink(name, 0)
				return link, nil
			}

			info.linkName = name
			name = link
			fi, err = os.Lstat(name)
			if err != nil {
				return "", nil
			}
		}

		
		
		for {
			info.wd, err = unix.Open(name, openMode, 0)
			if err == nil {
				break
			}
			if errors.Is(err, unix.EINTR) {
				continue
			}

			return "", err
		}

		info.isDir = fi.IsDir()
	}

	err := w.register([]int{info.wd}, unix.EV_ADD|unix.EV_CLEAR|unix.EV_ENABLE, flags)
	if err != nil {
		unix.Close(info.wd)
		return "", err
	}

	if !alreadyWatching {
		w.watches.add(name, info.linkName, info.wd, info.isDir)
	}

	
	
	if info.isDir {
		watchDir := (flags&unix.NOTE_WRITE) == unix.NOTE_WRITE &&
			(!alreadyWatching || (info.dirFlags&unix.NOTE_WRITE) != unix.NOTE_WRITE)
		w.watches.updateDirFlags(name, flags)

		if watchDir {
			if err := w.watchDirectoryFiles(name); err != nil {
				return "", err
			}
		}
	}
	return name, nil
}



func (w *kqueue) readEvents() {
	defer func() {
		close(w.Events)
		close(w.Errors)
		_ = unix.Close(w.kq)
		unix.Close(w.closepipe[0])
	}()

	eventBuffer := make([]unix.Kevent_t, 10)
	for {
		kevents, err := w.read(eventBuffer)
		
		if err != nil && err != unix.EINTR {
			if !w.sendError(fmt.Errorf("fsnotify.readEvents: %w", err)) {
				return
			}
		}

		for _, kevent := range kevents {
			var (
				wd   = int(kevent.Ident)
				mask = uint32(kevent.Fflags)
			)

			
			
			if wd == w.closepipe[0] {
				return
			}

			path, ok := w.watches.byWd(wd)
			if debug {
				internal.Debug(path.name, &kevent)
			}

			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			if !ok && kevent.Ident == 0 && runtime.GOOS == "darwin" {
				continue
			}

			event := w.newEvent(path.name, path.linkName, mask)

			if event.Has(Rename) || event.Has(Remove) {
				w.remove(event.Name, false)
				w.watches.markSeen(event.Name, false)
			}

			if path.isDir && event.Has(Write) && !event.Has(Remove) {
				w.dirChange(event.Name)
			} else if !w.sendEvent(event) {
				return
			}

			if event.Has(Remove) {
				
				
				if path.isDir {
					fileDir := filepath.Clean(event.Name)
					_, found := w.watches.byPath(fileDir)
					if found {
						
						
						
						
						
						
						
						
						
						
						
						
						
						
						
						err := w.dirChange(fileDir)
						if !w.sendError(err) {
							return
						}
					}
				} else {
					path := filepath.Clean(event.Name)
					if fi, err := os.Lstat(path); err == nil {
						err := w.sendCreateIfNew(path, fi)
						if !w.sendError(err) {
							return
						}
					}
				}
			}
		}
	}
}


func (w *kqueue) newEvent(name, linkName string, mask uint32) Event {
	e := Event{Name: name}
	if linkName != "" {
		
		
		e.Name = linkName
	}

	if mask&unix.NOTE_DELETE == unix.NOTE_DELETE {
		e.Op |= Remove
	}
	if mask&unix.NOTE_WRITE == unix.NOTE_WRITE {
		e.Op |= Write
	}
	if mask&unix.NOTE_RENAME == unix.NOTE_RENAME {
		e.Op |= Rename
	}
	if mask&unix.NOTE_ATTRIB == unix.NOTE_ATTRIB {
		e.Op |= Chmod
	}
	
	
	if e.Op.Has(Write) && e.Op.Has(Remove) {
		e.Op &^= Write
	}
	return e
}


func (w *kqueue) watchDirectoryFiles(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		path := filepath.Join(dirPath, f.Name())

		fi, err := f.Info()
		if err != nil {
			return fmt.Errorf("%q: %w", path, err)
		}

		cleanPath, err := w.internalWatch(path, fi)
		if err != nil {
			
			
			
			
			switch {
			case errors.Is(err, unix.EACCES) || errors.Is(err, unix.EPERM):
				cleanPath = filepath.Clean(path)
			default:
				return fmt.Errorf("%q: %w", path, err)
			}
		}

		w.watches.markSeen(cleanPath, true)
	}

	return nil
}





func (w *kqueue) dirChange(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		
		
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("fsnotify.dirChange: %w", err)
	}

	for _, f := range files {
		fi, err := f.Info()
		if err != nil {
			return fmt.Errorf("fsnotify.dirChange: %w", err)
		}

		err = w.sendCreateIfNew(filepath.Join(dir, fi.Name()), fi)
		if err != nil {
			
			if errors.Is(err, unix.EACCES) || errors.Is(err, unix.EPERM) {
				return nil
			}
			return fmt.Errorf("fsnotify.dirChange: %w", err)
		}
	}
	return nil
}



func (w *kqueue) sendCreateIfNew(path string, fi os.FileInfo) error {
	if !w.watches.seenBefore(path) {
		if !w.sendEvent(Event{Name: path, Op: Create}) {
			return nil
		}
	}

	
	path, err := w.internalWatch(path, fi)
	if err != nil {
		return err
	}
	w.watches.markSeen(path, true)
	return nil
}

func (w *kqueue) internalWatch(name string, fi os.FileInfo) (string, error) {
	if fi.IsDir() {
		
		
		info, _ := w.watches.byPath(name)
		return w.addWatch(name, info.dirFlags|unix.NOTE_DELETE|unix.NOTE_RENAME)
	}

	
	return w.addWatch(name, noteAllEvents)
}


func (w *kqueue) register(fds []int, flags int, fflags uint32) error {
	changes := make([]unix.Kevent_t, len(fds))
	for i, fd := range fds {
		
		unix.SetKevent(&changes[i], fd, unix.EVFILT_VNODE, flags)
		changes[i].Fflags = fflags
	}

	
	success, err := unix.Kevent(w.kq, changes, nil, nil)
	if success == -1 {
		return err
	}
	return nil
}


func (w *kqueue) read(events []unix.Kevent_t) ([]unix.Kevent_t, error) {
	n, err := unix.Kevent(w.kq, nil, events, nil)
	if err != nil {
		return nil, err
	}
	return events[0:n], nil
}

func (w *kqueue) xSupports(op Op) bool {
	if runtime.GOOS == "freebsd" {
		
	}
	if op.Has(xUnportableOpen) || op.Has(xUnportableRead) ||
		op.Has(xUnportableCloseWrite) || op.Has(xUnportableCloseRead) {
		return false
	}
	return true
}
