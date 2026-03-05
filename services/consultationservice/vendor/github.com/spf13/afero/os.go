













package afero

import (
	"os"
	"time"
)

var _ Lstater = (*OsFs)(nil)





type OsFs struct{}

func NewOsFs() Fs {
	return &OsFs{}
}

func (OsFs) Name() string { return "OsFs" }

func (OsFs) Create(name string) (File, error) {
	f, e := os.Create(name)
	if f == nil {
		
		
		return nil, e
	}
	return f, e
}

func (OsFs) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (OsFs) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (OsFs) Open(name string) (File, error) {
	f, e := os.Open(name)
	if f == nil {
		
		
		return nil, e
	}
	return f, e
}

func (OsFs) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	f, e := os.OpenFile(name, flag, perm)
	if f == nil {
		
		
		return nil, e
	}
	return f, e
}

func (OsFs) Remove(name string) error {
	return os.Remove(name)
}

func (OsFs) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (OsFs) Rename(oldname, newname string) error {
	return os.Rename(oldname, newname)
}

func (OsFs) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (OsFs) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}

func (OsFs) Chown(name string, uid, gid int) error {
	return os.Chown(name, uid, gid)
}

func (OsFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(name, atime, mtime)
}

func (OsFs) LstatIfPossible(name string) (os.FileInfo, bool, error) {
	fi, err := os.Lstat(name)
	return fi, true, err
}

func (OsFs) SymlinkIfPossible(oldname, newname string) error {
	return os.Symlink(oldname, newname)
}

func (OsFs) ReadlinkIfPossible(name string) (string, error) {
	return os.Readlink(name)
}
