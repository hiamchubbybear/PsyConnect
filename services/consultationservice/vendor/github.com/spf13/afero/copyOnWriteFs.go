package afero

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

var _ Lstater = (*CopyOnWriteFs)(nil)








type CopyOnWriteFs struct {
	base  Fs
	layer Fs
}

func NewCopyOnWriteFs(base Fs, layer Fs) Fs {
	return &CopyOnWriteFs{base: base, layer: layer}
}


func (u *CopyOnWriteFs) isBaseFile(name string) (bool, error) {
	if _, err := u.layer.Stat(name); err == nil {
		return false, nil
	}
	_, err := u.base.Stat(name)
	if err != nil {
		if oerr, ok := err.(*os.PathError); ok {
			if oerr.Err == os.ErrNotExist || oerr.Err == syscall.ENOENT || oerr.Err == syscall.ENOTDIR {
				return false, nil
			}
		}
		if err == syscall.ENOENT {
			return false, nil
		}
	}
	return true, err
}

func (u *CopyOnWriteFs) copyToLayer(name string) error {
	return copyToLayer(u.base, u.layer, name)
}

func (u *CopyOnWriteFs) Chtimes(name string, atime, mtime time.Time) error {
	b, err := u.isBaseFile(name)
	if err != nil {
		return err
	}
	if b {
		if err := u.copyToLayer(name); err != nil {
			return err
		}
	}
	return u.layer.Chtimes(name, atime, mtime)
}

func (u *CopyOnWriteFs) Chmod(name string, mode os.FileMode) error {
	b, err := u.isBaseFile(name)
	if err != nil {
		return err
	}
	if b {
		if err := u.copyToLayer(name); err != nil {
			return err
		}
	}
	return u.layer.Chmod(name, mode)
}

func (u *CopyOnWriteFs) Chown(name string, uid, gid int) error {
	b, err := u.isBaseFile(name)
	if err != nil {
		return err
	}
	if b {
		if err := u.copyToLayer(name); err != nil {
			return err
		}
	}
	return u.layer.Chown(name, uid, gid)
}

func (u *CopyOnWriteFs) Stat(name string) (os.FileInfo, error) {
	fi, err := u.layer.Stat(name)
	if err != nil {
		isNotExist := u.isNotExist(err)
		if isNotExist {
			return u.base.Stat(name)
		}
		return nil, err
	}
	return fi, nil
}

func (u *CopyOnWriteFs) LstatIfPossible(name string) (os.FileInfo, bool, error) {
	llayer, ok1 := u.layer.(Lstater)
	lbase, ok2 := u.base.(Lstater)

	if ok1 {
		fi, b, err := llayer.LstatIfPossible(name)
		if err == nil {
			return fi, b, nil
		}

		if !u.isNotExist(err) {
			return nil, b, err
		}
	}

	if ok2 {
		fi, b, err := lbase.LstatIfPossible(name)
		if err == nil {
			return fi, b, nil
		}
		if !u.isNotExist(err) {
			return nil, b, err
		}
	}

	fi, err := u.Stat(name)

	return fi, false, err
}

func (u *CopyOnWriteFs) SymlinkIfPossible(oldname, newname string) error {
	if slayer, ok := u.layer.(Linker); ok {
		return slayer.SymlinkIfPossible(oldname, newname)
	}

	return &os.LinkError{Op: "symlink", Old: oldname, New: newname, Err: ErrNoSymlink}
}

func (u *CopyOnWriteFs) ReadlinkIfPossible(name string) (string, error) {
	if rlayer, ok := u.layer.(LinkReader); ok {
		return rlayer.ReadlinkIfPossible(name)
	}

	if rbase, ok := u.base.(LinkReader); ok {
		return rbase.ReadlinkIfPossible(name)
	}

	return "", &os.PathError{Op: "readlink", Path: name, Err: ErrNoReadlink}
}

func (u *CopyOnWriteFs) isNotExist(err error) bool {
	if e, ok := err.(*os.PathError); ok {
		err = e.Err
	}
	if err == os.ErrNotExist || err == syscall.ENOENT || err == syscall.ENOTDIR {
		return true
	}
	return false
}


func (u *CopyOnWriteFs) Rename(oldname, newname string) error {
	b, err := u.isBaseFile(oldname)
	if err != nil {
		return err
	}
	if b {
		return syscall.EPERM
	}
	return u.layer.Rename(oldname, newname)
}




func (u *CopyOnWriteFs) Remove(name string) error {
	err := u.layer.Remove(name)
	switch err {
	case syscall.ENOENT:
		_, err = u.base.Stat(name)
		if err == nil {
			return syscall.EPERM
		}
		return syscall.ENOENT
	default:
		return err
	}
}

func (u *CopyOnWriteFs) RemoveAll(name string) error {
	err := u.layer.RemoveAll(name)
	switch err {
	case syscall.ENOENT:
		_, err = u.base.Stat(name)
		if err == nil {
			return syscall.EPERM
		}
		return syscall.ENOENT
	default:
		return err
	}
}

func (u *CopyOnWriteFs) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	b, err := u.isBaseFile(name)
	if err != nil {
		return nil, err
	}

	if flag&(os.O_WRONLY|os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_TRUNC) != 0 {
		if b {
			if err = u.copyToLayer(name); err != nil {
				return nil, err
			}
			return u.layer.OpenFile(name, flag, perm)
		}

		dir := filepath.Dir(name)
		isaDir, err := IsDir(u.base, dir)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		if isaDir {
			if err = u.layer.MkdirAll(dir, 0o777); err != nil {
				return nil, err
			}
			return u.layer.OpenFile(name, flag, perm)
		}

		isaDir, err = IsDir(u.layer, dir)
		if err != nil {
			return nil, err
		}
		if isaDir {
			return u.layer.OpenFile(name, flag, perm)
		}

		return nil, &os.PathError{Op: "open", Path: name, Err: syscall.ENOTDIR} 
	}
	if b {
		return u.base.OpenFile(name, flag, perm)
	}
	return u.layer.OpenFile(name, flag, perm)
}






func (u *CopyOnWriteFs) Open(name string) (File, error) {
	
	b, err := u.isBaseFile(name)
	if err != nil {
		return nil, err
	}

	
	if b {
		return u.base.Open(name)
	}

	
	dir, err := IsDir(u.layer, name)
	if err != nil {
		return nil, err
	}
	if !dir {
		return u.layer.Open(name)
	}

	
	
	
	

	
	dir, err = IsDir(u.base, name)
	if !dir || err != nil {
		return u.layer.Open(name)
	}

	
	
	bfile, bErr := u.base.Open(name)
	lfile, lErr := u.layer.Open(name)

	
	if bErr != nil || lErr != nil {
		return nil, fmt.Errorf("BaseErr: %v\nOverlayErr: %v", bErr, lErr)
	}

	return &UnionFile{Base: bfile, Layer: lfile}, nil
}

func (u *CopyOnWriteFs) Mkdir(name string, perm os.FileMode) error {
	dir, err := IsDir(u.base, name)
	if err != nil {
		return u.layer.MkdirAll(name, perm)
	}
	if dir {
		return ErrFileExists
	}
	return u.layer.MkdirAll(name, perm)
}

func (u *CopyOnWriteFs) Name() string {
	return "CopyOnWriteFs"
}

func (u *CopyOnWriteFs) MkdirAll(name string, perm os.FileMode) error {
	dir, err := IsDir(u.base, name)
	if err != nil {
		return u.layer.MkdirAll(name, perm)
	}
	if dir {
		
		return nil
	}
	return u.layer.MkdirAll(name, perm)
}

func (u *CopyOnWriteFs) Create(name string) (File, error) {
	return u.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o666)
}
