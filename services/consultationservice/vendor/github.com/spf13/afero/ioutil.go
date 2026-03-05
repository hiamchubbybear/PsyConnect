














package afero

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)


type byName []os.FileInfo

func (f byName) Len() int           { return len(f) }
func (f byName) Less(i, j int) bool { return f[i].Name() < f[j].Name() }
func (f byName) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }



func (a Afero) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ReadDir(a.Fs, dirname)
}

func ReadDir(fs Fs, dirname string) ([]os.FileInfo, error) {
	f, err := fs.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Sort(byName(list))
	return list, nil
}





func (a Afero) ReadFile(filename string) ([]byte, error) {
	return ReadFile(a.Fs, filename)
}

func ReadFile(fs Fs, filename string) ([]byte, error) {
	f, err := fs.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	
	
	var n int64

	if fi, err := f.Stat(); err == nil {
		
		if size := fi.Size(); size < 1e9 {
			n = size
		}
	}
	
	
	
	
	
	return readAll(f, n+bytes.MinRead)
}



func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, capacity))
	
	
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}





func ReadAll(r io.Reader) ([]byte, error) {
	return readAll(r, bytes.MinRead)
}




func (a Afero) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return WriteFile(a.Fs, filename, data, perm)
}

func WriteFile(fs Fs, filename string, data []byte, perm os.FileMode) error {
	f, err := fs.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}





var (
	randNum uint32
	randmu  sync.Mutex
)

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}

func nextRandom() string {
	randmu.Lock()
	r := randNum
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 
	randNum = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}












func (a Afero) TempFile(dir, pattern string) (f File, err error) {
	return TempFile(a.Fs, dir, pattern)
}

func TempFile(fs Fs, dir, pattern string) (f File, err error) {
	if dir == "" {
		dir = os.TempDir()
	}

	var prefix, suffix string
	if pos := strings.LastIndex(pattern, "*"); pos != -1 {
		prefix, suffix = pattern[:pos], pattern[pos+1:]
	} else {
		prefix = pattern
	}

	nconflict := 0
	for i := 0; i < 10000; i++ {
		name := filepath.Join(dir, prefix+nextRandom()+suffix)
		f, err = fs.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o600)
		if os.IsExist(err) {
			if nconflict++; nconflict > 10 {
				randmu.Lock()
				randNum = reseed()
				randmu.Unlock()
			}
			continue
		}
		break
	}
	return
}








func (a Afero) TempDir(dir, prefix string) (name string, err error) {
	return TempDir(a.Fs, dir, prefix)
}

func TempDir(fs Fs, dir, prefix string) (name string, err error) {
	if dir == "" {
		dir = os.TempDir()
	}

	nconflict := 0
	for i := 0; i < 10000; i++ {
		try := filepath.Join(dir, prefix+nextRandom())
		err = fs.Mkdir(try, 0o700)
		if os.IsExist(err) {
			if nconflict++; nconflict > 10 {
				randmu.Lock()
				randNum = reseed()
				randmu.Unlock()
			}
			continue
		}
		if err == nil {
			name = try
		}
		break
	}
	return
}
