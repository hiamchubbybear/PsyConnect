












package afero

import (
	"errors"
)








type Symlinker interface {
	Lstater
	Linker
	LinkReader
}





type Linker interface {
	SymlinkIfPossible(oldname, newname string) error
}




var ErrNoSymlink = errors.New("symlink not supported")



type LinkReader interface {
	ReadlinkIfPossible(name string) (string, error)
}




var ErrNoReadlink = errors.New("readlink not supported")
