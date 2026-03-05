












package afero

import (
	"os"
)






type Lstater interface {
	LstatIfPossible(name string) (os.FileInfo, bool, error)
}
