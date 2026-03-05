package afero

import (
	"io"
	"os"
	"path/filepath"
	"syscall"
)














type UnionFile struct {
	Base   File
	Layer  File
	Merger DirsMerger
	off    int
	files  []os.FileInfo
}

func (f *UnionFile) Close() error {
	
	
	
	if f.Base != nil {
		f.Base.Close()
	}
	if f.Layer != nil {
		return f.Layer.Close()
	}
	return BADFD
}

func (f *UnionFile) Read(s []byte) (int, error) {
	if f.Layer != nil {
		n, err := f.Layer.Read(s)
		if (err == nil || err == io.EOF) && f.Base != nil {
			
			
			if _, seekErr := f.Base.Seek(int64(n), io.SeekCurrent); seekErr != nil {
				
				
				err = seekErr
			}
		}
		return n, err
	}
	if f.Base != nil {
		return f.Base.Read(s)
	}
	return 0, BADFD
}

func (f *UnionFile) ReadAt(s []byte, o int64) (int, error) {
	if f.Layer != nil {
		n, err := f.Layer.ReadAt(s, o)
		if (err == nil || err == io.EOF) && f.Base != nil {
			_, err = f.Base.Seek(o+int64(n), io.SeekStart)
		}
		return n, err
	}
	if f.Base != nil {
		return f.Base.ReadAt(s, o)
	}
	return 0, BADFD
}

func (f *UnionFile) Seek(o int64, w int) (pos int64, err error) {
	if f.Layer != nil {
		pos, err = f.Layer.Seek(o, w)
		if (err == nil || err == io.EOF) && f.Base != nil {
			_, err = f.Base.Seek(o, w)
		}
		return pos, err
	}
	if f.Base != nil {
		return f.Base.Seek(o, w)
	}
	return 0, BADFD
}

func (f *UnionFile) Write(s []byte) (n int, err error) {
	if f.Layer != nil {
		n, err = f.Layer.Write(s)
		if err == nil && f.Base != nil { 
			_, err = f.Base.Write(s)
		}
		return n, err
	}
	if f.Base != nil {
		return f.Base.Write(s)
	}
	return 0, BADFD
}

func (f *UnionFile) WriteAt(s []byte, o int64) (n int, err error) {
	if f.Layer != nil {
		n, err = f.Layer.WriteAt(s, o)
		if err == nil && f.Base != nil {
			_, err = f.Base.WriteAt(s, o)
		}
		return n, err
	}
	if f.Base != nil {
		return f.Base.WriteAt(s, o)
	}
	return 0, BADFD
}

func (f *UnionFile) Name() string {
	if f.Layer != nil {
		return f.Layer.Name()
	}
	return f.Base.Name()
}




type DirsMerger func(lofi, bofi []os.FileInfo) ([]os.FileInfo, error)

var defaultUnionMergeDirsFn = func(lofi, bofi []os.FileInfo) ([]os.FileInfo, error) {
	files := make(map[string]os.FileInfo)

	for _, fi := range lofi {
		files[fi.Name()] = fi
	}

	for _, fi := range bofi {
		if _, exists := files[fi.Name()]; !exists {
			files[fi.Name()] = fi
		}
	}

	rfi := make([]os.FileInfo, len(files))

	i := 0
	for _, fi := range files {
		rfi[i] = fi
		i++
	}

	return rfi, nil
}




func (f *UnionFile) Readdir(c int) (ofi []os.FileInfo, err error) {
	var merge DirsMerger = f.Merger
	if merge == nil {
		merge = defaultUnionMergeDirsFn
	}

	if f.off == 0 {
		var lfi []os.FileInfo
		if f.Layer != nil {
			lfi, err = f.Layer.Readdir(-1)
			if err != nil {
				return nil, err
			}
		}

		var bfi []os.FileInfo
		if f.Base != nil {
			bfi, err = f.Base.Readdir(-1)
			if err != nil {
				return nil, err
			}

		}
		merged, err := merge(lfi, bfi)
		if err != nil {
			return nil, err
		}
		f.files = append(f.files, merged...)
	}
	files := f.files[f.off:]

	if c <= 0 {
		return files, nil
	}

	if len(files) == 0 {
		return nil, io.EOF
	}

	if c > len(files) {
		c = len(files)
	}

	defer func() { f.off += c }()
	return files[:c], nil
}

func (f *UnionFile) Readdirnames(c int) ([]string, error) {
	rfi, err := f.Readdir(c)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, fi := range rfi {
		names = append(names, fi.Name())
	}
	return names, nil
}

func (f *UnionFile) Stat() (os.FileInfo, error) {
	if f.Layer != nil {
		return f.Layer.Stat()
	}
	if f.Base != nil {
		return f.Base.Stat()
	}
	return nil, BADFD
}

func (f *UnionFile) Sync() (err error) {
	if f.Layer != nil {
		err = f.Layer.Sync()
		if err == nil && f.Base != nil {
			err = f.Base.Sync()
		}
		return err
	}
	if f.Base != nil {
		return f.Base.Sync()
	}
	return BADFD
}

func (f *UnionFile) Truncate(s int64) (err error) {
	if f.Layer != nil {
		err = f.Layer.Truncate(s)
		if err == nil && f.Base != nil {
			err = f.Base.Truncate(s)
		}
		return err
	}
	if f.Base != nil {
		return f.Base.Truncate(s)
	}
	return BADFD
}

func (f *UnionFile) WriteString(s string) (n int, err error) {
	if f.Layer != nil {
		n, err = f.Layer.WriteString(s)
		if err == nil && f.Base != nil {
			_, err = f.Base.WriteString(s)
		}
		return n, err
	}
	if f.Base != nil {
		return f.Base.WriteString(s)
	}
	return 0, BADFD
}

func copyFile(base Fs, layer Fs, name string, bfh File) error {
	
	exists, err := Exists(layer, filepath.Dir(name))
	if err != nil {
		return err
	}
	if !exists {
		err = layer.MkdirAll(filepath.Dir(name), 0o777) 
		if err != nil {
			return err
		}
	}

	
	lfh, err := layer.Create(name)
	if err != nil {
		return err
	}
	n, err := io.Copy(lfh, bfh)
	if err != nil {
		
		layer.Remove(name)
		lfh.Close()
		return err
	}

	bfi, err := bfh.Stat()
	if err != nil || bfi.Size() != n {
		layer.Remove(name)
		lfh.Close()
		return syscall.EIO
	}

	err = lfh.Close()
	if err != nil {
		layer.Remove(name)
		lfh.Close()
		return err
	}
	return layer.Chtimes(name, bfi.ModTime(), bfi.ModTime())
}

func copyToLayer(base Fs, layer Fs, name string) error {
	bfh, err := base.Open(name)
	if err != nil {
		return err
	}
	defer bfh.Close()

	return copyFile(base, layer, name, bfh)
}

func copyFileToLayer(base Fs, layer Fs, name string, flag int, perm os.FileMode) error {
	bfh, err := base.OpenFile(name, flag, perm)
	if err != nil {
		return err
	}
	defer bfh.Close()

	return copyFile(base, layer, name, bfh)
}
