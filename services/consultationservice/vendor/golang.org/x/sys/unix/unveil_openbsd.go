



package unix

import "fmt"





func Unveil(path string, flags string) error {
	if err := supportsUnveil(); err != nil {
		return err
	}
	pathPtr, err := BytePtrFromString(path)
	if err != nil {
		return err
	}
	flagsPtr, err := BytePtrFromString(flags)
	if err != nil {
		return err
	}
	return unveil(pathPtr, flagsPtr)
}



func UnveilBlock() error {
	if err := supportsUnveil(); err != nil {
		return err
	}
	return unveil(nil, nil)
}



func supportsUnveil() error {
	maj, min, err := majmin()
	if err != nil {
		return err
	}

	
	if maj < 6 || (maj == 6 && min <= 3) {
		return fmt.Errorf("cannot call Unveil on OpenBSD %d.%d", maj, min)
	}

	return nil
}
